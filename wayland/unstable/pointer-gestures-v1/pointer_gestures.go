// Generated by go-wayland-scanner
// https://github.com/rajveermalviya/go-wayland/cmd/go-wayland-scanner
// XML file : https://raw.githubusercontent.com/wayland-project/wayland-protocols/1.27/unstable/pointer-gestures/pointer-gestures-unstable-v1.xml
//
// pointer_gestures_unstable_v1 Protocol Copyright:

package pointer_gestures

import "github.com/rajveermalviya/go-wayland/wayland/client"

// PointerGestures : touchpad gestures
//
// A global interface to provide semantic touchpad gestures for a given
// pointer.
//
// Three gestures are currently supported: swipe, pinch, and hold.
// Pinch and swipe gestures follow a three-stage cycle: begin, update,
// end, hold gestures follow a two-stage cycle: begin and end. All
// gestures are identified by a unique id.
//
// Warning! The protocol described in this file is experimental and
// backward incompatible changes may be made. Backward compatible changes
// may be added together with the corresponding interface version bump.
// Backward incompatible changes are done by bumping the version number in
// the protocol and interface names and resetting the interface version.
// Once the protocol is to be declared stable, the 'z' prefix and the
// version number in the protocol and interface names are removed and the
// interface version number is reset.
type PointerGestures struct {
	client.BaseProxy
}

// NewPointerGestures : touchpad gestures
//
// A global interface to provide semantic touchpad gestures for a given
// pointer.
//
// Three gestures are currently supported: swipe, pinch, and hold.
// Pinch and swipe gestures follow a three-stage cycle: begin, update,
// end, hold gestures follow a two-stage cycle: begin and end. All
// gestures are identified by a unique id.
//
// Warning! The protocol described in this file is experimental and
// backward incompatible changes may be made. Backward compatible changes
// may be added together with the corresponding interface version bump.
// Backward incompatible changes are done by bumping the version number in
// the protocol and interface names and resetting the interface version.
// Once the protocol is to be declared stable, the 'z' prefix and the
// version number in the protocol and interface names are removed and the
// interface version number is reset.
func NewPointerGestures(ctx *client.Context) *PointerGestures {
	zwpPointerGesturesV1 := &PointerGestures{}
	ctx.Register(zwpPointerGesturesV1)
	return zwpPointerGesturesV1
}

// GetSwipeGesture : get swipe gesture
//
// Create a swipe gesture object. See the
// wl_pointer_gesture_swipe interface for details.
func (i *PointerGestures) GetSwipeGesture(pointer *client.Pointer) (*PointerGestureSwipe, error) {
	id := NewPointerGestureSwipe(i.Context())
	const opcode = 0
	const _reqBufLen = 8 + 4 + 4
	var _reqBuf [_reqBufLen]byte
	l := 0
	client.PutUint32(_reqBuf[l:4], i.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], uint32(_reqBufLen<<16|opcode&0x0000ffff))
	l += 4
	client.PutUint32(_reqBuf[l:l+4], id.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], pointer.ID())
	l += 4
	err := i.Context().WriteMsg(_reqBuf[:], nil)
	return id, err
}

// GetPinchGesture : get pinch gesture
//
// Create a pinch gesture object. See the
// wl_pointer_gesture_pinch interface for details.
func (i *PointerGestures) GetPinchGesture(pointer *client.Pointer) (*PointerGesturePinch, error) {
	id := NewPointerGesturePinch(i.Context())
	const opcode = 1
	const _reqBufLen = 8 + 4 + 4
	var _reqBuf [_reqBufLen]byte
	l := 0
	client.PutUint32(_reqBuf[l:4], i.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], uint32(_reqBufLen<<16|opcode&0x0000ffff))
	l += 4
	client.PutUint32(_reqBuf[l:l+4], id.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], pointer.ID())
	l += 4
	err := i.Context().WriteMsg(_reqBuf[:], nil)
	return id, err
}

// Release : destroy the pointer gesture object
//
// Destroy the pointer gesture object. Swipe, pinch and hold objects
// created via this gesture object remain valid.
func (i *PointerGestures) Release() error {
	defer i.Context().Unregister(i)
	const opcode = 2
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

// GetHoldGesture : get hold gesture
//
// Create a hold gesture object. See the
// wl_pointer_gesture_hold interface for details.
func (i *PointerGestures) GetHoldGesture(pointer *client.Pointer) (*PointerGestureHold, error) {
	id := NewPointerGestureHold(i.Context())
	const opcode = 3
	const _reqBufLen = 8 + 4 + 4
	var _reqBuf [_reqBufLen]byte
	l := 0
	client.PutUint32(_reqBuf[l:4], i.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], uint32(_reqBufLen<<16|opcode&0x0000ffff))
	l += 4
	client.PutUint32(_reqBuf[l:l+4], id.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], pointer.ID())
	l += 4
	err := i.Context().WriteMsg(_reqBuf[:], nil)
	return id, err
}

// PointerGestureSwipe : a swipe gesture object
//
// A swipe gesture object notifies a client about a multi-finger swipe
// gesture detected on an indirect input device such as a touchpad.
// The gesture is usually initiated by multiple fingers moving in the
// same direction but once initiated the direction may change.
// The precise conditions of when such a gesture is detected are
// implementation-dependent.
//
// A gesture consists of three stages: begin, update (optional) and end.
// There cannot be multiple simultaneous hold, pinch or swipe gestures on a
// same pointer/seat, how compositors prevent these situations is
// implementation-dependent.
//
// A gesture may be cancelled by the compositor or the hardware.
// Clients should not consider performing permanent or irreversible
// actions until the end of a gesture has been received.
type PointerGestureSwipe struct {
	client.BaseProxy
	beginHandlers  []PointerGestureSwipeBeginHandlerFunc
	updateHandlers []PointerGestureSwipeUpdateHandlerFunc
	endHandlers    []PointerGestureSwipeEndHandlerFunc
}

// NewPointerGestureSwipe : a swipe gesture object
//
// A swipe gesture object notifies a client about a multi-finger swipe
// gesture detected on an indirect input device such as a touchpad.
// The gesture is usually initiated by multiple fingers moving in the
// same direction but once initiated the direction may change.
// The precise conditions of when such a gesture is detected are
// implementation-dependent.
//
// A gesture consists of three stages: begin, update (optional) and end.
// There cannot be multiple simultaneous hold, pinch or swipe gestures on a
// same pointer/seat, how compositors prevent these situations is
// implementation-dependent.
//
// A gesture may be cancelled by the compositor or the hardware.
// Clients should not consider performing permanent or irreversible
// actions until the end of a gesture has been received.
func NewPointerGestureSwipe(ctx *client.Context) *PointerGestureSwipe {
	zwpPointerGestureSwipeV1 := &PointerGestureSwipe{}
	ctx.Register(zwpPointerGestureSwipeV1)
	return zwpPointerGestureSwipeV1
}

// Destroy : destroy the pointer swipe gesture object
func (i *PointerGestureSwipe) Destroy() error {
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

// PointerGestureSwipeBeginEvent : multi-finger swipe begin
//
// This event is sent when a multi-finger swipe gesture is detected
// on the device.
type PointerGestureSwipeBeginEvent struct {
	Serial  uint32
	Time    uint32
	Surface *client.Surface
	Fingers uint32
}
type PointerGestureSwipeBeginHandlerFunc func(PointerGestureSwipeBeginEvent)

// AddBeginHandler : adds handler for PointerGestureSwipeBeginEvent
func (i *PointerGestureSwipe) AddBeginHandler(f PointerGestureSwipeBeginHandlerFunc) {
	if f == nil {
		return
	}

	i.beginHandlers = append(i.beginHandlers, f)
}

// PointerGestureSwipeUpdateEvent : multi-finger swipe motion
//
// This event is sent when a multi-finger swipe gesture changes the
// position of the logical center.
//
// The dx and dy coordinates are relative coordinates of the logical
// center of the gesture compared to the previous event.
type PointerGestureSwipeUpdateEvent struct {
	Time uint32
	Dx   float64
	Dy   float64
}
type PointerGestureSwipeUpdateHandlerFunc func(PointerGestureSwipeUpdateEvent)

// AddUpdateHandler : adds handler for PointerGestureSwipeUpdateEvent
func (i *PointerGestureSwipe) AddUpdateHandler(f PointerGestureSwipeUpdateHandlerFunc) {
	if f == nil {
		return
	}

	i.updateHandlers = append(i.updateHandlers, f)
}

// PointerGestureSwipeEndEvent : multi-finger swipe end
//
// This event is sent when a multi-finger swipe gesture ceases to
// be valid. This may happen when one or more fingers are lifted or
// the gesture is cancelled.
//
// When a gesture is cancelled, the client should undo state changes
// caused by this gesture. What causes a gesture to be cancelled is
// implementation-dependent.
type PointerGestureSwipeEndEvent struct {
	Serial    uint32
	Time      uint32
	Cancelled int32
}
type PointerGestureSwipeEndHandlerFunc func(PointerGestureSwipeEndEvent)

// AddEndHandler : adds handler for PointerGestureSwipeEndEvent
func (i *PointerGestureSwipe) AddEndHandler(f PointerGestureSwipeEndHandlerFunc) {
	if f == nil {
		return
	}

	i.endHandlers = append(i.endHandlers, f)
}

func (i *PointerGestureSwipe) Dispatch(opcode uint16, fd uintptr, data []byte) {
	switch opcode {
	case 0:
		if len(i.beginHandlers) == 0 {
			return
		}
		var e PointerGestureSwipeBeginEvent
		l := 0
		e.Serial = client.Uint32(data[l : l+4])
		l += 4
		e.Time = client.Uint32(data[l : l+4])
		l += 4
		e.Surface = i.Context().GetProxy(client.Uint32(data[l : l+4])).(*client.Surface)
		l += 4
		e.Fingers = client.Uint32(data[l : l+4])
		l += 4
		for _, f := range i.beginHandlers {
			f(e)
		}
	case 1:
		if len(i.updateHandlers) == 0 {
			return
		}
		var e PointerGestureSwipeUpdateEvent
		l := 0
		e.Time = client.Uint32(data[l : l+4])
		l += 4
		e.Dx = client.Fixed(data[l : l+4])
		l += 4
		e.Dy = client.Fixed(data[l : l+4])
		l += 4
		for _, f := range i.updateHandlers {
			f(e)
		}
	case 2:
		if len(i.endHandlers) == 0 {
			return
		}
		var e PointerGestureSwipeEndEvent
		l := 0
		e.Serial = client.Uint32(data[l : l+4])
		l += 4
		e.Time = client.Uint32(data[l : l+4])
		l += 4
		e.Cancelled = int32(client.Uint32(data[l : l+4]))
		l += 4
		for _, f := range i.endHandlers {
			f(e)
		}
	}
}

// PointerGesturePinch : a pinch gesture object
//
// A pinch gesture object notifies a client about a multi-finger pinch
// gesture detected on an indirect input device such as a touchpad.
// The gesture is usually initiated by multiple fingers moving towards
// each other or away from each other, or by two or more fingers rotating
// around a logical center of gravity. The precise conditions of when
// such a gesture is detected are implementation-dependent.
//
// A gesture consists of three stages: begin, update (optional) and end.
// There cannot be multiple simultaneous hold, pinch or swipe gestures on a
// same pointer/seat, how compositors prevent these situations is
// implementation-dependent.
//
// A gesture may be cancelled by the compositor or the hardware.
// Clients should not consider performing permanent or irreversible
// actions until the end of a gesture has been received.
type PointerGesturePinch struct {
	client.BaseProxy
	beginHandlers  []PointerGesturePinchBeginHandlerFunc
	updateHandlers []PointerGesturePinchUpdateHandlerFunc
	endHandlers    []PointerGesturePinchEndHandlerFunc
}

// NewPointerGesturePinch : a pinch gesture object
//
// A pinch gesture object notifies a client about a multi-finger pinch
// gesture detected on an indirect input device such as a touchpad.
// The gesture is usually initiated by multiple fingers moving towards
// each other or away from each other, or by two or more fingers rotating
// around a logical center of gravity. The precise conditions of when
// such a gesture is detected are implementation-dependent.
//
// A gesture consists of three stages: begin, update (optional) and end.
// There cannot be multiple simultaneous hold, pinch or swipe gestures on a
// same pointer/seat, how compositors prevent these situations is
// implementation-dependent.
//
// A gesture may be cancelled by the compositor or the hardware.
// Clients should not consider performing permanent or irreversible
// actions until the end of a gesture has been received.
func NewPointerGesturePinch(ctx *client.Context) *PointerGesturePinch {
	zwpPointerGesturePinchV1 := &PointerGesturePinch{}
	ctx.Register(zwpPointerGesturePinchV1)
	return zwpPointerGesturePinchV1
}

// Destroy : destroy the pinch gesture object
func (i *PointerGesturePinch) Destroy() error {
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

// PointerGesturePinchBeginEvent : multi-finger pinch begin
//
// This event is sent when a multi-finger pinch gesture is detected
// on the device.
type PointerGesturePinchBeginEvent struct {
	Serial  uint32
	Time    uint32
	Surface *client.Surface
	Fingers uint32
}
type PointerGesturePinchBeginHandlerFunc func(PointerGesturePinchBeginEvent)

// AddBeginHandler : adds handler for PointerGesturePinchBeginEvent
func (i *PointerGesturePinch) AddBeginHandler(f PointerGesturePinchBeginHandlerFunc) {
	if f == nil {
		return
	}

	i.beginHandlers = append(i.beginHandlers, f)
}

// PointerGesturePinchUpdateEvent : multi-finger pinch motion
//
// This event is sent when a multi-finger pinch gesture changes the
// position of the logical center, the rotation or the relative scale.
//
// The dx and dy coordinates are relative coordinates in the
// surface coordinate space of the logical center of the gesture.
//
// The scale factor is an absolute scale compared to the
// pointer_gesture_pinch.begin event, e.g. a scale of 2 means the fingers
// are now twice as far apart as on pointer_gesture_pinch.begin.
//
// The rotation is the relative angle in degrees clockwise compared to the previous
// pointer_gesture_pinch.begin or pointer_gesture_pinch.update event.
type PointerGesturePinchUpdateEvent struct {
	Time     uint32
	Dx       float64
	Dy       float64
	Scale    float64
	Rotation float64
}
type PointerGesturePinchUpdateHandlerFunc func(PointerGesturePinchUpdateEvent)

// AddUpdateHandler : adds handler for PointerGesturePinchUpdateEvent
func (i *PointerGesturePinch) AddUpdateHandler(f PointerGesturePinchUpdateHandlerFunc) {
	if f == nil {
		return
	}

	i.updateHandlers = append(i.updateHandlers, f)
}

// PointerGesturePinchEndEvent : multi-finger pinch end
//
// This event is sent when a multi-finger pinch gesture ceases to
// be valid. This may happen when one or more fingers are lifted or
// the gesture is cancelled.
//
// When a gesture is cancelled, the client should undo state changes
// caused by this gesture. What causes a gesture to be cancelled is
// implementation-dependent.
type PointerGesturePinchEndEvent struct {
	Serial    uint32
	Time      uint32
	Cancelled int32
}
type PointerGesturePinchEndHandlerFunc func(PointerGesturePinchEndEvent)

// AddEndHandler : adds handler for PointerGesturePinchEndEvent
func (i *PointerGesturePinch) AddEndHandler(f PointerGesturePinchEndHandlerFunc) {
	if f == nil {
		return
	}

	i.endHandlers = append(i.endHandlers, f)
}

func (i *PointerGesturePinch) Dispatch(opcode uint16, fd uintptr, data []byte) {
	switch opcode {
	case 0:
		if len(i.beginHandlers) == 0 {
			return
		}
		var e PointerGesturePinchBeginEvent
		l := 0
		e.Serial = client.Uint32(data[l : l+4])
		l += 4
		e.Time = client.Uint32(data[l : l+4])
		l += 4
		e.Surface = i.Context().GetProxy(client.Uint32(data[l : l+4])).(*client.Surface)
		l += 4
		e.Fingers = client.Uint32(data[l : l+4])
		l += 4
		for _, f := range i.beginHandlers {
			f(e)
		}
	case 1:
		if len(i.updateHandlers) == 0 {
			return
		}
		var e PointerGesturePinchUpdateEvent
		l := 0
		e.Time = client.Uint32(data[l : l+4])
		l += 4
		e.Dx = client.Fixed(data[l : l+4])
		l += 4
		e.Dy = client.Fixed(data[l : l+4])
		l += 4
		e.Scale = client.Fixed(data[l : l+4])
		l += 4
		e.Rotation = client.Fixed(data[l : l+4])
		l += 4
		for _, f := range i.updateHandlers {
			f(e)
		}
	case 2:
		if len(i.endHandlers) == 0 {
			return
		}
		var e PointerGesturePinchEndEvent
		l := 0
		e.Serial = client.Uint32(data[l : l+4])
		l += 4
		e.Time = client.Uint32(data[l : l+4])
		l += 4
		e.Cancelled = int32(client.Uint32(data[l : l+4]))
		l += 4
		for _, f := range i.endHandlers {
			f(e)
		}
	}
}

// PointerGestureHold : a hold gesture object
//
// A hold gesture object notifies a client about a single- or
// multi-finger hold gesture detected on an indirect input device such as
// a touchpad. The gesture is usually initiated by one or more fingers
// being held down without significant movement. The precise conditions
// of when such a gesture is detected are implementation-dependent.
//
// In particular, this gesture may be used to cancel kinetic scrolling.
//
// A hold gesture consists of two stages: begin and end. Unlike pinch and
// swipe there is no update stage.
// There cannot be multiple simultaneous hold, pinch or swipe gestures on a
// same pointer/seat, how compositors prevent these situations is
// implementation-dependent.
//
// A gesture may be cancelled by the compositor or the hardware.
// Clients should not consider performing permanent or irreversible
// actions until the end of a gesture has been received.
type PointerGestureHold struct {
	client.BaseProxy
	beginHandlers []PointerGestureHoldBeginHandlerFunc
	endHandlers   []PointerGestureHoldEndHandlerFunc
}

// NewPointerGestureHold : a hold gesture object
//
// A hold gesture object notifies a client about a single- or
// multi-finger hold gesture detected on an indirect input device such as
// a touchpad. The gesture is usually initiated by one or more fingers
// being held down without significant movement. The precise conditions
// of when such a gesture is detected are implementation-dependent.
//
// In particular, this gesture may be used to cancel kinetic scrolling.
//
// A hold gesture consists of two stages: begin and end. Unlike pinch and
// swipe there is no update stage.
// There cannot be multiple simultaneous hold, pinch or swipe gestures on a
// same pointer/seat, how compositors prevent these situations is
// implementation-dependent.
//
// A gesture may be cancelled by the compositor or the hardware.
// Clients should not consider performing permanent or irreversible
// actions until the end of a gesture has been received.
func NewPointerGestureHold(ctx *client.Context) *PointerGestureHold {
	zwpPointerGestureHoldV1 := &PointerGestureHold{}
	ctx.Register(zwpPointerGestureHoldV1)
	return zwpPointerGestureHoldV1
}

// Destroy : destroy the hold gesture object
func (i *PointerGestureHold) Destroy() error {
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

// PointerGestureHoldBeginEvent : multi-finger hold begin
//
// This event is sent when a hold gesture is detected on the device.
type PointerGestureHoldBeginEvent struct {
	Serial  uint32
	Time    uint32
	Surface *client.Surface
	Fingers uint32
}
type PointerGestureHoldBeginHandlerFunc func(PointerGestureHoldBeginEvent)

// AddBeginHandler : adds handler for PointerGestureHoldBeginEvent
func (i *PointerGestureHold) AddBeginHandler(f PointerGestureHoldBeginHandlerFunc) {
	if f == nil {
		return
	}

	i.beginHandlers = append(i.beginHandlers, f)
}

// PointerGestureHoldEndEvent : multi-finger hold end
//
// This event is sent when a hold gesture ceases to
// be valid. This may happen when the holding fingers are lifted or
// the gesture is cancelled, for example if the fingers move past an
// implementation-defined threshold, the finger count changes or the hold
// gesture changes into a different type of gesture.
//
// When a gesture is cancelled, the client may need to undo state changes
// caused by this gesture. What causes a gesture to be cancelled is
// implementation-dependent.
type PointerGestureHoldEndEvent struct {
	Serial    uint32
	Time      uint32
	Cancelled int32
}
type PointerGestureHoldEndHandlerFunc func(PointerGestureHoldEndEvent)

// AddEndHandler : adds handler for PointerGestureHoldEndEvent
func (i *PointerGestureHold) AddEndHandler(f PointerGestureHoldEndHandlerFunc) {
	if f == nil {
		return
	}

	i.endHandlers = append(i.endHandlers, f)
}

func (i *PointerGestureHold) Dispatch(opcode uint16, fd uintptr, data []byte) {
	switch opcode {
	case 0:
		if len(i.beginHandlers) == 0 {
			return
		}
		var e PointerGestureHoldBeginEvent
		l := 0
		e.Serial = client.Uint32(data[l : l+4])
		l += 4
		e.Time = client.Uint32(data[l : l+4])
		l += 4
		e.Surface = i.Context().GetProxy(client.Uint32(data[l : l+4])).(*client.Surface)
		l += 4
		e.Fingers = client.Uint32(data[l : l+4])
		l += 4
		for _, f := range i.beginHandlers {
			f(e)
		}
	case 1:
		if len(i.endHandlers) == 0 {
			return
		}
		var e PointerGestureHoldEndEvent
		l := 0
		e.Serial = client.Uint32(data[l : l+4])
		l += 4
		e.Time = client.Uint32(data[l : l+4])
		l += 4
		e.Cancelled = int32(client.Uint32(data[l : l+4]))
		l += 4
		for _, f := range i.endHandlers {
			f(e)
		}
	}
}
