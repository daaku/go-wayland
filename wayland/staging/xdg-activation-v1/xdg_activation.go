// Generated by go-wayland-scanner
// https://github.com/rajveermalviya/go-wayland/cmd/go-wayland-scanner
// XML file : https://gitlab.freedesktop.org/wayland/wayland-protocols/-/raw/1.36/staging/xdg-activation/xdg-activation-v1.xml
//
// xdg_activation_v1 Protocol Copyright:
//
// Copyright © 2020 Aleix Pol Gonzalez <aleixpol@kde.org>
// Copyright © 2020 Carlos Garnacho <carlosg@gnome.org>
//
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice (including the next
// paragraph) shall be included in all copies or substantial portions of the
// Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL
// THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package xdg_activation

import "github.com/rajveermalviya/go-wayland/wayland/client"

// Activation : interface for activating surfaces
//
// A global interface used for informing the compositor about applications
// being activated or started, or for applications to request to be
// activated.
type Activation struct {
	client.BaseProxy
}

// NewActivation : interface for activating surfaces
//
// A global interface used for informing the compositor about applications
// being activated or started, or for applications to request to be
// activated.
func NewActivation(ctx *client.Context) *Activation {
	xdgActivationV1 := &Activation{}
	ctx.Register(xdgActivationV1)
	return xdgActivationV1
}

// Destroy : destroy the xdg_activation object
//
// Notify the compositor that the xdg_activation object will no longer be
// used.
//
// The child objects created via this interface are unaffected and should
// be destroyed separately.
func (i *Activation) Destroy() error {
	defer i.Context().Unregister(i)
	const opcode = 0
	const _reqBufLen = 8
	var _reqBuf [_reqBufLen]byte
	l := 0
	client.PutUint32(_reqBuf[l:4], i.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], uint32(_reqBufLen<<16|opcode&0x0000ffff))
	l += 4
	err := i.Context().WriteMsg(_reqBuf[:], nil)
	return err
}

// GetActivationToken : requests a token
//
// Creates an xdg_activation_token_v1 object that will provide
// the initiating client with a unique token for this activation. This
// token should be offered to the clients to be activated.
func (i *Activation) GetActivationToken() (*ActivationToken, error) {
	id := NewActivationToken(i.Context())
	const opcode = 1
	const _reqBufLen = 8 + 4
	var _reqBuf [_reqBufLen]byte
	l := 0
	client.PutUint32(_reqBuf[l:4], i.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], uint32(_reqBufLen<<16|opcode&0x0000ffff))
	l += 4
	client.PutUint32(_reqBuf[l:l+4], id.ID())
	l += 4
	err := i.Context().WriteMsg(_reqBuf[:], nil)
	return id, err
}

// Activate : notify new interaction being available
//
// Requests surface activation. It's up to the compositor to display
// this information as desired, for example by placing the surface above
// the rest.
//
// The compositor may know who requested this by checking the activation
// token and might decide not to follow through with the activation if it's
// considered unwanted.
//
// Compositors can ignore unknown activation tokens when an invalid
// token is passed.
//
//	token: the activation token of the initiating client
//	surface: the wl_surface to activate
func (i *Activation) Activate(token string, surface *client.Surface) error {
	const opcode = 2
	tokenLen := client.PaddedLen(len(token) + 1)
	_reqBufLen := 8 + (4 + tokenLen) + 4
	_reqBuf := make([]byte, _reqBufLen)
	l := 0
	client.PutUint32(_reqBuf[l:4], i.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], uint32(_reqBufLen<<16|opcode&0x0000ffff))
	l += 4
	client.PutString(_reqBuf[l:l+(4+tokenLen)], token, tokenLen)
	l += (4 + tokenLen)
	client.PutUint32(_reqBuf[l:l+4], surface.ID())
	l += 4
	err := i.Context().WriteMsg(_reqBuf, nil)
	return err
}

// ActivationToken : an exported activation handle
//
// An object for setting up a token and receiving a token handle that can
// be passed as an activation token to another client.
//
// The object is created using the xdg_activation_v1.get_activation_token
// request. This object should then be populated with the app_id, surface
// and serial information and committed. The compositor shall then issue a
// done event with the token. In case the request's parameters are invalid,
// the compositor will provide an invalid token.
type ActivationToken struct {
	client.BaseProxy
	doneHandler ActivationTokenDoneHandlerFunc
}

// NewActivationToken : an exported activation handle
//
// An object for setting up a token and receiving a token handle that can
// be passed as an activation token to another client.
//
// The object is created using the xdg_activation_v1.get_activation_token
// request. This object should then be populated with the app_id, surface
// and serial information and committed. The compositor shall then issue a
// done event with the token. In case the request's parameters are invalid,
// the compositor will provide an invalid token.
func NewActivationToken(ctx *client.Context) *ActivationToken {
	xdgActivationTokenV1 := &ActivationToken{}
	ctx.Register(xdgActivationTokenV1)
	return xdgActivationTokenV1
}

// SetSerial : specifies the seat and serial of the activating event
//
// Provides information about the seat and serial event that requested the
// token.
//
// The serial can come from an input or focus event. For instance, if a
// click triggers the launch of a third-party client, the launcher client
// should send a set_serial request with the serial and seat from the
// wl_pointer.button event.
//
// Some compositors might refuse to activate toplevels when the token
// doesn't have a valid and recent enough event serial.
//
// Must be sent before commit. This information is optional.
//
//	serial: the serial of the event that triggered the activation
//	seat: the wl_seat of the event
func (i *ActivationToken) SetSerial(serial uint32, seat *client.Seat) error {
	const opcode = 0
	const _reqBufLen = 8 + 4 + 4
	var _reqBuf [_reqBufLen]byte
	l := 0
	client.PutUint32(_reqBuf[l:4], i.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], uint32(_reqBufLen<<16|opcode&0x0000ffff))
	l += 4
	client.PutUint32(_reqBuf[l:l+4], uint32(serial))
	l += 4
	client.PutUint32(_reqBuf[l:l+4], seat.ID())
	l += 4
	err := i.Context().WriteMsg(_reqBuf[:], nil)
	return err
}

// SetAppId : specifies the application being activated
//
// The requesting client can specify an app_id to associate the token
// being created with it.
//
// Must be sent before commit. This information is optional.
//
//	appId: the application id of the client being activated.
func (i *ActivationToken) SetAppId(appId string) error {
	const opcode = 1
	appIdLen := client.PaddedLen(len(appId) + 1)
	_reqBufLen := 8 + (4 + appIdLen)
	_reqBuf := make([]byte, _reqBufLen)
	l := 0
	client.PutUint32(_reqBuf[l:4], i.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], uint32(_reqBufLen<<16|opcode&0x0000ffff))
	l += 4
	client.PutString(_reqBuf[l:l+(4+appIdLen)], appId, appIdLen)
	l += (4 + appIdLen)
	err := i.Context().WriteMsg(_reqBuf, nil)
	return err
}

// SetSurface : specifies the surface requesting activation
//
// This request sets the surface requesting the activation. Note, this is
// different from the surface that will be activated.
//
// Some compositors might refuse to activate toplevels when the token
// doesn't have a requesting surface.
//
// Must be sent before commit. This information is optional.
//
//	surface: the requesting surface
func (i *ActivationToken) SetSurface(surface *client.Surface) error {
	const opcode = 2
	const _reqBufLen = 8 + 4
	var _reqBuf [_reqBufLen]byte
	l := 0
	client.PutUint32(_reqBuf[l:4], i.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], uint32(_reqBufLen<<16|opcode&0x0000ffff))
	l += 4
	client.PutUint32(_reqBuf[l:l+4], surface.ID())
	l += 4
	err := i.Context().WriteMsg(_reqBuf[:], nil)
	return err
}

// Commit : issues the token request
//
// Requests an activation token based on the different parameters that
// have been offered through set_serial, set_surface and set_app_id.
func (i *ActivationToken) Commit() error {
	const opcode = 3
	const _reqBufLen = 8
	var _reqBuf [_reqBufLen]byte
	l := 0
	client.PutUint32(_reqBuf[l:4], i.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], uint32(_reqBufLen<<16|opcode&0x0000ffff))
	l += 4
	err := i.Context().WriteMsg(_reqBuf[:], nil)
	return err
}

// Destroy : destroy the xdg_activation_token_v1 object
//
// Notify the compositor that the xdg_activation_token_v1 object will no
// longer be used. The received token stays valid.
func (i *ActivationToken) Destroy() error {
	defer i.Context().Unregister(i)
	const opcode = 4
	const _reqBufLen = 8
	var _reqBuf [_reqBufLen]byte
	l := 0
	client.PutUint32(_reqBuf[l:4], i.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], uint32(_reqBufLen<<16|opcode&0x0000ffff))
	l += 4
	err := i.Context().WriteMsg(_reqBuf[:], nil)
	return err
}

type ActivationTokenError uint32

// ActivationTokenError :
const (
	// ActivationTokenErrorAlreadyUsed : The token has already been used previously
	ActivationTokenErrorAlreadyUsed ActivationTokenError = 0
)

func (e ActivationTokenError) Name() string {
	switch e {
	case ActivationTokenErrorAlreadyUsed:
		return "already_used"
	default:
		return ""
	}
}

func (e ActivationTokenError) Value() string {
	switch e {
	case ActivationTokenErrorAlreadyUsed:
		return "0"
	default:
		return ""
	}
}

func (e ActivationTokenError) String() string {
	return e.Name() + "=" + e.Value()
}

// ActivationTokenDoneEvent : the exported activation token
//
// The 'done' event contains the unique token of this activation request
// and notifies that the provider is done.
type ActivationTokenDoneEvent struct {
	Token string
}
type ActivationTokenDoneHandlerFunc func(ActivationTokenDoneEvent)

// SetDoneHandler : sets handler for ActivationTokenDoneEvent
func (i *ActivationToken) SetDoneHandler(f ActivationTokenDoneHandlerFunc) {
	i.doneHandler = f
}

func (i *ActivationToken) Dispatch(opcode uint32, fd int, data []byte) {
	switch opcode {
	case 0:
		if i.doneHandler == nil {
			return
		}
		var e ActivationTokenDoneEvent
		l := 0
		tokenLen := client.PaddedLen(int(client.Uint32(data[l : l+4])))
		l += 4
		e.Token = client.String(data[l : l+tokenLen])
		l += tokenLen

		i.doneHandler(e)
	}
}
