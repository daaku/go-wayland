package idle_inhibit

//go:generate go run github.com/rajveermalviya/go-wayland/cmd/go-wayland-scanner -pkg idle_inhibit -prefix zwp -suffix v1 -o idle_inhibit.go -i https://gitlab.freedesktop.org/wayland/wayland-protocols/-/raw/1.36/unstable/idle-inhibit/idle-inhibit-unstable-v1.xml
