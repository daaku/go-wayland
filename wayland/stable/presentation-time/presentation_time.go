// Generated by go-wayland-scanner
// https://github.com/rajveermalviya/go-wayland/cmd/go-wayland-scanner
// XML file : https://raw.githubusercontent.com/wayland-project/wayland-protocols/1.27/stable/presentation-time/presentation-time.xml
//
// presentation_time Protocol Copyright:
//
// Copyright © 2013-2014 Collabora, Ltd.
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

package presentation_time

import "github.com/rajveermalviya/go-wayland/wayland/client"

// Presentation : timed presentation related wl_surface requests
//
// The main feature of this interface is accurate presentation
// timing feedback to ensure smooth video playback while maintaining
// audio/video synchronization. Some features use the concept of a
// presentation clock, which is defined in the
// presentation.clock_id event.
//
// A content update for a wl_surface is submitted by a
// wl_surface.commit request. Request 'feedback' associates with
// the wl_surface.commit and provides feedback on the content
// update, particularly the final realized presentation time.
//
// When the final realized presentation time is available, e.g.
// after a framebuffer flip completes, the requested
// presentation_feedback.presented events are sent. The final
// presentation time can differ from the compositor's predicted
// display update time and the update's target time, especially
// when the compositor misses its target vertical blanking period.
type Presentation struct {
	client.BaseProxy
	clockIdHandlers []PresentationClockIdHandlerFunc
}

// NewPresentation : timed presentation related wl_surface requests
//
// The main feature of this interface is accurate presentation
// timing feedback to ensure smooth video playback while maintaining
// audio/video synchronization. Some features use the concept of a
// presentation clock, which is defined in the
// presentation.clock_id event.
//
// A content update for a wl_surface is submitted by a
// wl_surface.commit request. Request 'feedback' associates with
// the wl_surface.commit and provides feedback on the content
// update, particularly the final realized presentation time.
//
// When the final realized presentation time is available, e.g.
// after a framebuffer flip completes, the requested
// presentation_feedback.presented events are sent. The final
// presentation time can differ from the compositor's predicted
// display update time and the update's target time, especially
// when the compositor misses its target vertical blanking period.
func NewPresentation(ctx *client.Context) *Presentation {
	wpPresentation := &Presentation{}
	ctx.Register(wpPresentation)
	return wpPresentation
}

// Destroy : unbind from the presentation interface
//
// Informs the server that the client will no longer be using
// this protocol object. Existing objects created by this object
// are not affected.
func (i *Presentation) Destroy() error {
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

// Feedback : request presentation feedback information
//
// Request presentation feedback for the current content submission
// on the given surface. This creates a new presentation_feedback
// object, which will deliver the feedback information once. If
// multiple presentation_feedback objects are created for the same
// submission, they will all deliver the same information.
//
// For details on what information is returned, see the
// presentation_feedback interface.
//
//	surface: target surface
func (i *Presentation) Feedback(surface *client.Surface) (*PresentationFeedback, error) {
	callback := NewPresentationFeedback(i.Context())
	const opcode = 1
	const _reqBufLen = 8 + 4 + 4
	var _reqBuf [_reqBufLen]byte
	l := 0
	client.PutUint32(_reqBuf[l:4], i.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], uint32(_reqBufLen<<16|opcode&0x0000ffff))
	l += 4
	client.PutUint32(_reqBuf[l:l+4], surface.ID())
	l += 4
	client.PutUint32(_reqBuf[l:l+4], callback.ID())
	l += 4
	err := i.Context().WriteMsg(_reqBuf[:], nil)
	return callback, err
}

type PresentationError uint32

// PresentationError : fatal presentation errors
//
// These fatal protocol errors may be emitted in response to
// illegal presentation requests.
const (
	// PresentationErrorInvalidTimestamp : invalid value in tv_nsec
	PresentationErrorInvalidTimestamp PresentationError = 0
	// PresentationErrorInvalidFlag : invalid flag
	PresentationErrorInvalidFlag PresentationError = 1
)

func (e PresentationError) Name() string {
	switch e {
	case PresentationErrorInvalidTimestamp:
		return "invalid_timestamp"
	case PresentationErrorInvalidFlag:
		return "invalid_flag"
	default:
		return ""
	}
}

func (e PresentationError) Value() string {
	switch e {
	case PresentationErrorInvalidTimestamp:
		return "0"
	case PresentationErrorInvalidFlag:
		return "1"
	default:
		return ""
	}
}

func (e PresentationError) String() string {
	return e.Name() + "=" + e.Value()
}

// PresentationClockIdEvent : clock ID for timestamps
//
// This event tells the client in which clock domain the
// compositor interprets the timestamps used by the presentation
// extension. This clock is called the presentation clock.
//
// The compositor sends this event when the client binds to the
// presentation interface. The presentation clock does not change
// during the lifetime of the client connection.
//
// The clock identifier is platform dependent. On Linux/glibc,
// the identifier value is one of the clockid_t values accepted
// by clock_gettime(). clock_gettime() is defined by
// POSIX.1-2001.
//
// Timestamps in this clock domain are expressed as tv_sec_hi,
// tv_sec_lo, tv_nsec triples, each component being an unsigned
// 32-bit value. Whole seconds are in tv_sec which is a 64-bit
// value combined from tv_sec_hi and tv_sec_lo, and the
// additional fractional part in tv_nsec as nanoseconds. Hence,
// for valid timestamps tv_nsec must be in [0, 999999999].
//
// Note that clock_id applies only to the presentation clock,
// and implies nothing about e.g. the timestamps used in the
// Wayland core protocol input events.
//
// Compositors should prefer a clock which does not jump and is
// not slewed e.g. by NTP. The absolute value of the clock is
// irrelevant. Precision of one millisecond or better is
// recommended. Clients must be able to query the current clock
// value directly, not by asking the compositor.
type PresentationClockIdEvent struct {
	ClkId uint32
}
type PresentationClockIdHandlerFunc func(PresentationClockIdEvent)

// AddClockIdHandler : adds handler for PresentationClockIdEvent
func (i *Presentation) AddClockIdHandler(f PresentationClockIdHandlerFunc) {
	if f == nil {
		return
	}

	i.clockIdHandlers = append(i.clockIdHandlers, f)
}

func (i *Presentation) Dispatch(opcode uint16, fd uintptr, data []byte) {
	switch opcode {
	case 0:
		if len(i.clockIdHandlers) == 0 {
			return
		}
		var e PresentationClockIdEvent
		l := 0
		e.ClkId = client.Uint32(data[l : l+4])
		l += 4
		for _, f := range i.clockIdHandlers {
			f(e)
		}
	}
}

// PresentationFeedback : presentation time feedback event
//
// A presentation_feedback object returns an indication that a
// wl_surface content update has become visible to the user.
// One object corresponds to one content update submission
// (wl_surface.commit). There are two possible outcomes: the
// content update is presented to the user, and a presentation
// timestamp delivered; or, the user did not see the content
// update because it was superseded or its surface destroyed,
// and the content update is discarded.
//
// Once a presentation_feedback object has delivered a 'presented'
// or 'discarded' event it is automatically destroyed.
type PresentationFeedback struct {
	client.BaseProxy
	syncOutputHandlers []PresentationFeedbackSyncOutputHandlerFunc
	presentedHandlers  []PresentationFeedbackPresentedHandlerFunc
	discardedHandlers  []PresentationFeedbackDiscardedHandlerFunc
}

// NewPresentationFeedback : presentation time feedback event
//
// A presentation_feedback object returns an indication that a
// wl_surface content update has become visible to the user.
// One object corresponds to one content update submission
// (wl_surface.commit). There are two possible outcomes: the
// content update is presented to the user, and a presentation
// timestamp delivered; or, the user did not see the content
// update because it was superseded or its surface destroyed,
// and the content update is discarded.
//
// Once a presentation_feedback object has delivered a 'presented'
// or 'discarded' event it is automatically destroyed.
func NewPresentationFeedback(ctx *client.Context) *PresentationFeedback {
	wpPresentationFeedback := &PresentationFeedback{}
	ctx.Register(wpPresentationFeedback)
	return wpPresentationFeedback
}

func (i *PresentationFeedback) Destroy() error {
	i.Context().Unregister(i)
	return nil
}

type PresentationFeedbackKind uint32

// PresentationFeedbackKind : bitmask of flags in presented event
//
// These flags provide information about how the presentation of
// the related content update was done. The intent is to help
// clients assess the reliability of the feedback and the visual
// quality with respect to possible tearing and timings.
const (
	PresentationFeedbackKindVsync        PresentationFeedbackKind = 0x1
	PresentationFeedbackKindHwClock      PresentationFeedbackKind = 0x2
	PresentationFeedbackKindHwCompletion PresentationFeedbackKind = 0x4
	PresentationFeedbackKindZeroCopy     PresentationFeedbackKind = 0x8
)

func (e PresentationFeedbackKind) Name() string {
	switch e {
	case PresentationFeedbackKindVsync:
		return "vsync"
	case PresentationFeedbackKindHwClock:
		return "hw_clock"
	case PresentationFeedbackKindHwCompletion:
		return "hw_completion"
	case PresentationFeedbackKindZeroCopy:
		return "zero_copy"
	default:
		return ""
	}
}

func (e PresentationFeedbackKind) Value() string {
	switch e {
	case PresentationFeedbackKindVsync:
		return "0x1"
	case PresentationFeedbackKindHwClock:
		return "0x2"
	case PresentationFeedbackKindHwCompletion:
		return "0x4"
	case PresentationFeedbackKindZeroCopy:
		return "0x8"
	default:
		return ""
	}
}

func (e PresentationFeedbackKind) String() string {
	return e.Name() + "=" + e.Value()
}

// PresentationFeedbackSyncOutputEvent : presentation synchronized to this output
//
// As presentation can be synchronized to only one output at a
// time, this event tells which output it was. This event is only
// sent prior to the presented event.
//
// As clients may bind to the same global wl_output multiple
// times, this event is sent for each bound instance that matches
// the synchronized output. If a client has not bound to the
// right wl_output global at all, this event is not sent.
type PresentationFeedbackSyncOutputEvent struct {
	Output *client.Output
}
type PresentationFeedbackSyncOutputHandlerFunc func(PresentationFeedbackSyncOutputEvent)

// AddSyncOutputHandler : adds handler for PresentationFeedbackSyncOutputEvent
func (i *PresentationFeedback) AddSyncOutputHandler(f PresentationFeedbackSyncOutputHandlerFunc) {
	if f == nil {
		return
	}

	i.syncOutputHandlers = append(i.syncOutputHandlers, f)
}

// PresentationFeedbackPresentedEvent : the content update was displayed
//
// The associated content update was displayed to the user at the
// indicated time (tv_sec_hi/lo, tv_nsec). For the interpretation of
// the timestamp, see presentation.clock_id event.
//
// The timestamp corresponds to the time when the content update
// turned into light the first time on the surface's main output.
// Compositors may approximate this from the framebuffer flip
// completion events from the system, and the latency of the
// physical display path if known.
//
// This event is preceded by all related sync_output events
// telling which output's refresh cycle the feedback corresponds
// to, i.e. the main output for the surface. Compositors are
// recommended to choose the output containing the largest part
// of the wl_surface, or keeping the output they previously
// chose. Having a stable presentation output association helps
// clients predict future output refreshes (vblank).
//
// The 'refresh' argument gives the compositor's prediction of how
// many nanoseconds after tv_sec, tv_nsec the very next output
// refresh may occur. This is to further aid clients in
// predicting future refreshes, i.e., estimating the timestamps
// targeting the next few vblanks. If such prediction cannot
// usefully be done, the argument is zero.
//
// If the output does not have a constant refresh rate, explicit
// video mode switches excluded, then the refresh argument must
// be zero.
//
// The 64-bit value combined from seq_hi and seq_lo is the value
// of the output's vertical retrace counter when the content
// update was first scanned out to the display. This value must
// be compatible with the definition of MSC in
// GLX_OML_sync_control specification. Note, that if the display
// path has a non-zero latency, the time instant specified by
// this counter may differ from the timestamp's.
//
// If the output does not have a concept of vertical retrace or a
// refresh cycle, or the output device is self-refreshing without
// a way to query the refresh count, then the arguments seq_hi
// and seq_lo must be zero.
type PresentationFeedbackPresentedEvent struct {
	TvSecHi uint32
	TvSecLo uint32
	TvNsec  uint32
	Refresh uint32
	SeqHi   uint32
	SeqLo   uint32
	Flags   uint32
}
type PresentationFeedbackPresentedHandlerFunc func(PresentationFeedbackPresentedEvent)

// AddPresentedHandler : adds handler for PresentationFeedbackPresentedEvent
func (i *PresentationFeedback) AddPresentedHandler(f PresentationFeedbackPresentedHandlerFunc) {
	if f == nil {
		return
	}

	i.presentedHandlers = append(i.presentedHandlers, f)
}

// PresentationFeedbackDiscardedEvent : the content update was not displayed
//
// The content update was never displayed to the user.
type PresentationFeedbackDiscardedEvent struct{}
type PresentationFeedbackDiscardedHandlerFunc func(PresentationFeedbackDiscardedEvent)

// AddDiscardedHandler : adds handler for PresentationFeedbackDiscardedEvent
func (i *PresentationFeedback) AddDiscardedHandler(f PresentationFeedbackDiscardedHandlerFunc) {
	if f == nil {
		return
	}

	i.discardedHandlers = append(i.discardedHandlers, f)
}

func (i *PresentationFeedback) Dispatch(opcode uint16, fd uintptr, data []byte) {
	switch opcode {
	case 0:
		if len(i.syncOutputHandlers) == 0 {
			return
		}
		var e PresentationFeedbackSyncOutputEvent
		l := 0
		e.Output = i.Context().GetProxy(client.Uint32(data[l : l+4])).(*client.Output)
		l += 4
		for _, f := range i.syncOutputHandlers {
			f(e)
		}
	case 1:
		if len(i.presentedHandlers) == 0 {
			return
		}
		var e PresentationFeedbackPresentedEvent
		l := 0
		e.TvSecHi = client.Uint32(data[l : l+4])
		l += 4
		e.TvSecLo = client.Uint32(data[l : l+4])
		l += 4
		e.TvNsec = client.Uint32(data[l : l+4])
		l += 4
		e.Refresh = client.Uint32(data[l : l+4])
		l += 4
		e.SeqHi = client.Uint32(data[l : l+4])
		l += 4
		e.SeqLo = client.Uint32(data[l : l+4])
		l += 4
		e.Flags = client.Uint32(data[l : l+4])
		l += 4
		for _, f := range i.presentedHandlers {
			f(e)
		}
	case 2:
		if len(i.discardedHandlers) == 0 {
			return
		}
		var e PresentationFeedbackDiscardedEvent
		for _, f := range i.discardedHandlers {
			f(e)
		}
	}
}
