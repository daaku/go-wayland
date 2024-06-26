package main

import (
	"fmt"
	"image"
	"log"
	"os"

	"github.com/daaku/swizzle"
	"github.com/fogleman/gg"
	"github.com/pkg/errors"
	"github.com/rajveermalviya/go-wayland/wayland/client"
	xdg_shell "github.com/rajveermalviya/go-wayland/wayland/stable/xdg-shell"
	"golang.org/x/sys/unix"
)

// Global app state
type appState struct {
	appID         string
	title         string
	width, height int32
	frame         *image.RGBA
	exit          bool
	skipDraw      bool

	display    *client.Display
	registry   *client.Registry
	shm        *client.Shm
	compositor *client.Compositor
	xdgWmBase  *xdg_shell.WmBase
	seat       *client.Seat

	surface     *client.Surface
	xdgSurface  *xdg_shell.Surface
	xdgTopLevel *xdg_shell.Toplevel

	keyboard *client.Keyboard
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func run() error {
	const S = 1000
	dc := gg.NewContext(S, S)
	dc.DrawCircle(500, 500, 400)
	dc.SetRGB(128, 0, 0)
	dc.Fill()
	dc.SetRGB(255, 255, 255)
	if err := dc.LoadFontFace("/usr/share/fonts/TTF/Roboto-Black.ttf", 96); err != nil {
		return errors.WithStack(err)
	}
	dc.DrawStringAnchored("Hello, world!", S/2, S/2, 0.5, 0.5)
	frameImage := dc.Image().(*image.RGBA)
	frameRect := frameImage.Bounds()

	app := &appState{
		title:  "gg",
		appID:  "osd_imageviewer",
		width:  int32(frameRect.Dx()),
		height: int32(frameRect.Dy()),
		frame:  frameImage,
	}

	if err := app.initWindow(); err != nil {
		return err
	}

	// Start the dispatch loop
	for !app.exit {
		app.dispatch()
	}

	app.cleanup()
	return nil
}

func (app *appState) initWindow() error {
	display, err := client.Connect("")
	if err != nil {
		return errors.Errorf("unable to connect to wayland server: %v", err)
	}
	app.display = display
	display.SetErrorHandler(app.HandleDisplayError)

	registry, err := app.display.GetRegistry()
	if err != nil {
		return errors.WithStack(err)
	}
	app.registry = registry

	registry.SetGlobalHandler(app.HandleRegistryGlobal)
	app.displayRoundTrip() // Wait for interfaces to register
	app.displayRoundTrip() // Wait for handler events

	// Create a wl_surface for toplevel window
	surface, err := app.compositor.CreateSurface()
	if err != nil {
		return errors.WithStack(err)
	}
	app.surface = surface

	// attach wl_surface to xdg_wmbase to get toplevel handle
	xdgSurface, err := app.xdgWmBase.GetXdgSurface(surface)
	if err != nil {
		return errors.WithStack(err)
	}
	app.xdgSurface = xdgSurface
	xdgSurface.SetConfigureHandler(app.HandleSurfaceConfigure)

	xdgTopLevel, err := xdgSurface.GetToplevel()
	if err != nil {
		return errors.WithStack(err)
	}
	app.xdgTopLevel = xdgTopLevel
	xdgTopLevel.SetConfigureHandler(app.HandleToplevelConfigure) // for window resize
	xdgTopLevel.SetCloseHandler(app.HandleToplevelClose)

	if err := xdgTopLevel.SetTitle(app.title); err != nil {
		return errors.WithStack(err)
	}
	if err := xdgTopLevel.SetAppId(app.appID); err != nil {
		return errors.WithStack(err)
	}
	if err := app.surface.Commit(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (app *appState) dispatch() {
	app.display.Context().Dispatch()
}

func (app *appState) context() *client.Context {
	return app.display.Context()
}

func (app *appState) HandleRegistryGlobal(e client.RegistryGlobalEvent) {
	switch e.Interface {
	case "wl_compositor":
		compositor := client.NewCompositor(app.context())
		err := app.registry.Bind(e.Name, e.Interface, e.Version, compositor)
		if err != nil {
			log.Fatalf("unable to bind wl_compositor interface: %v", err)
		}
		app.compositor = compositor
	case "wl_shm":
		shm := client.NewShm(app.context())
		err := app.registry.Bind(e.Name, e.Interface, e.Version, shm)
		if err != nil {
			log.Fatalf("unable to bind wl_shm interface: %v", err)
		}
		app.shm = shm
	case "xdg_wm_base":
		xdgWmBase := xdg_shell.NewWmBase(app.context())
		err := app.registry.Bind(e.Name, e.Interface, e.Version, xdgWmBase)
		if err != nil {
			log.Fatalf("unable to bind xdg_wm_base interface: %v", err)
		}
		app.xdgWmBase = xdgWmBase
		xdgWmBase.SetPingHandler(app.HandleWmBasePing)
	case "wl_seat":
		seat := client.NewSeat(app.context())
		err := app.registry.Bind(e.Name, e.Interface, e.Version, seat)
		if err != nil {
			log.Fatalf("unable to bind wl_seat interface: %v", err)
		}
		app.seat = seat
		seat.SetCapabilitiesHandler(app.HandleSeatCapabilities)
	}
}

func (app *appState) HandleSurfaceConfigure(e xdg_shell.SurfaceConfigureEvent) {
	if err := app.xdgSurface.AckConfigure(e.Serial); err != nil {
		log.Fatal("unable to ack xdg surface configure")
	}

	if app.skipDraw {
		return
	}

	buffer := app.drawFrame()
	if err := app.surface.Attach(buffer, 0, 0); err != nil {
		log.Fatalf("unable to attach buffer to surface: %v", err)
	}
	if err := app.surface.Commit(); err != nil {
		log.Fatalf("unable to commit surface state: %v", err)
	}
}

func (app *appState) HandleToplevelConfigure(e xdg_shell.ToplevelConfigureEvent) {
	width := e.Width
	height := e.Height

	if width == 0 || height == 0 {
		// Compositor is deferring to us
		return
	}

	if width == app.width && height == app.height {
		// No need to resize
		return
	}

	// Update app size
	app.width = width
	app.height = height
}

func (app *appState) drawFrame() *client.Buffer {
	app.skipDraw = true

	stride := app.width * 4
	size := stride * app.height

	file, err := anonTempFile(int64(size))
	if err != nil {
		log.Fatalf("unable to create a temporary file: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			logPrintf("unable to close file: %v", err)
		}
	}()

	data, err := unix.Mmap(int(file.Fd()), 0, int(size), unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		log.Fatalf("unable to create mapping: %v", err)
	}
	defer func() {
		if err := unix.Munmap(data); err != nil {
			logPrintf("unable to delete mapping: %v", err)
		}
	}()

	pool, err := app.shm.CreatePool(int(file.Fd()), size)
	if err != nil {
		log.Fatalf("unable to create shm pool: %v", err)
	}
	defer func() {
		if err := pool.Destroy(); err != nil {
			logPrintf("unable to destroy shm pool: %v", err)
		}
	}()

	buf, err := pool.CreateBuffer(0, app.width, app.height, stride, uint32(client.ShmFormatArgb8888))
	if err != nil {
		log.Fatalf("unable to create client.Buffer from shm pool: %v", err)
	}

	// Convert RGBA to BGRA
	copy(data, app.frame.Pix)
	swizzle.BGRA(data)

	buf.SetReleaseHandler(func(client.BufferReleaseEvent) {
		if err := buf.Destroy(); err != nil {
			logPrintf("unable to destroy buffer: %v", err)
		}
	})

	return buf
}

func (app *appState) HandleSeatCapabilities(e client.SeatCapabilitiesEvent) {
	haveKeyboard := (e.Capabilities * uint32(client.SeatCapabilityKeyboard)) != 0
	if haveKeyboard && app.keyboard == nil {
		app.attachKeyboard()
	} else if !haveKeyboard && app.keyboard != nil {
		app.releaseKeyboard()
	}
}

// HandleDisplayError handles client.Display errors
func (*appState) HandleDisplayError(e client.DisplayErrorEvent) {
	// Just log.Fatal for now
	log.Fatalf("display error event: %v", e)
}

func (app *appState) HandleWmBasePing(e xdg_shell.WmBasePingEvent) {
	app.xdgWmBase.Pong(e.Serial)
}

func (app *appState) HandleToplevelClose(xdg_shell.ToplevelCloseEvent) {
	app.exit = true
}

func (app *appState) displayRoundTrip() {
	callback, err := app.display.Sync()
	if err != nil {
		log.Fatalf("unable to get sync callback: %v", err)
	}
	defer func() {
		if err := callback.Destroy(); err != nil {
			logPrintln("unable to destroy callback:", err)
		}
	}()

	done := false
	callback.SetDoneHandler(func(client.CallbackDoneEvent) {
		done = true
	})

	// Wait for callback to return
	for !done {
		app.dispatch()
	}
}

func (app *appState) cleanup() {
	if app.keyboard != nil {
		app.releaseKeyboard()
	}

	if app.xdgTopLevel != nil {
		if err := app.xdgTopLevel.Destroy(); err != nil {
			logPrintln("unable to destroy xdg_toplevel:", err)
		}
		app.xdgTopLevel = nil
	}

	if app.xdgSurface != nil {
		if err := app.xdgSurface.Destroy(); err != nil {
			logPrintln("unable to destroy xdg_surface:", err)
		}
		app.xdgSurface = nil
	}

	if app.surface != nil {
		if err := app.surface.Destroy(); err != nil {
			logPrintln("unable to destroy wl_surface:", err)
		}
		app.surface = nil
	}

	// Release wl_seat handlers
	if app.seat != nil {
		if err := app.seat.Release(); err != nil {
			logPrintln("unable to destroy wl_seat:", err)
		}
		app.seat = nil
	}

	// Release xdg_wmbase
	if app.xdgWmBase != nil {
		if err := app.xdgWmBase.Destroy(); err != nil {
			logPrintln("unable to destroy xdg_wm_base:", err)
		}
		app.xdgWmBase = nil
	}

	if app.shm != nil {
		if err := app.shm.Destroy(); err != nil {
			logPrintln("unable to destroy wl_shm:", err)
		}
		app.shm = nil
	}

	if app.compositor != nil {
		if err := app.compositor.Destroy(); err != nil {
			logPrintln("unable to destroy wl_compositor:", err)
		}
		app.compositor = nil
	}

	if app.registry != nil {
		if err := app.registry.Destroy(); err != nil {
			logPrintln("unable to destroy wl_registry:", err)
		}
		app.registry = nil
	}

	if app.display != nil {
		if err := app.display.Destroy(); err != nil {
			logPrintln("unable to destroy wl_display:", err)
		}
	}

	if err := app.context().Close(); err != nil {
		logPrintln("unable to close wayland context:", err)
	}
}

func (app *appState) attachKeyboard() {
	keyboard, err := app.seat.GetKeyboard()
	if err != nil {
		log.Fatal("unable to register keyboard interface")
	}
	app.keyboard = keyboard
	keyboard.SetKeyHandler(app.HandleKeyboardKey)
}

func (app *appState) releaseKeyboard() {
	if err := app.keyboard.Release(); err != nil {
		logPrintln("unable to release keyboard interface")
	}
	app.keyboard = nil
}

func (app *appState) HandleKeyboardKey(e client.KeyboardKeyEvent) {
	// close on "esc"
	if e.Key == 1 {
		app.exit = true
	}
}

func anonTempFile(size int64) (*os.File, error) {
	dir := os.Getenv("XDG_RUNTIME_DIR")
	if dir == "" {
		return nil, errors.New("XDG_RUNTIME_DIR is not defined in env")
	}
	file, err := os.CreateTemp(dir, "wl_shm_go_*")
	if err != nil {
		return nil, err
	}
	err = file.Truncate(size)
	if err != nil {
		return nil, err
	}
	err = os.Remove(file.Name())
	if err != nil {
		return nil, err
	}
	return file, nil
}

var logDisabled = os.Getenv("LOG_DISABLED") == "1"

func logPrintln(v ...interface{}) {
	if !logDisabled {
		log.Println(v...)
	}
}

func logPrintf(format string, v ...interface{}) {
	if !logDisabled {
		log.Printf(format, v...)
	}
}
