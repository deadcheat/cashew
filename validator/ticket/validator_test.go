package ticket

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/deadcheat/cashew/mocks/timer"
	globaltimer "github.com/deadcheat/cashew/timer"
	"github.com/golang/mock/gomock"

	"github.com/deadcheat/cashew"
	"github.com/deadcheat/cashew/errors"
	"github.com/deadcheat/cashew/foundation"
)

func TestValidateTicketIsNil(t *testing.T) {

	testee := New()

	err := testee.Validate(nil)
	if err != errors.ErrTicketNotFound {
		t.Errorf("Validate() returned unexpected error %+v", err)
	}
}

func TestValidateLoginTicket(t *testing.T) {
	loadSetting()
	now := time.Now()
	c := gomock.NewController(t)
	tm := timer.NewMockTimeWrapper(c)

	tm.EXPECT().Now().AnyTimes().Return(now)

	globaltimer.Local = tm

	testTicket := &cashew.Ticket{
		Type:      cashew.TicketTypeLogin,
		CreatedAt: now.Add(time.Duration(-1*foundation.App().LoginTicketExpire) * time.Second),
	}
	testee := New()

	err := testee.Validate(testTicket)
	if err != nil {
		t.Errorf("Validate() returned unexpected error %+v", err)
	}

	testTicketForError := &cashew.Ticket{
		Type:      cashew.TicketTypeLogin,
		CreatedAt: now.Add(time.Duration(-1*foundation.App().LoginTicketExpire-1) * time.Second),
	}

	err = testee.Validate(testTicketForError)
	if err != errors.ErrTicketHasBeenExpired {
		t.Errorf("Validate() returned unexpected error %+v", err)
	}
}

func TestValidatePGT(t *testing.T) {
	loadSetting()
	now := time.Now()
	c := gomock.NewController(t)
	tm := timer.NewMockTimeWrapper(c)

	tm.EXPECT().Now().AnyTimes().Return(now)

	globaltimer.Local = tm

	testTime := now.Add(time.Duration(-1*foundation.App().GrantingDefaultExpire) * time.Second)

	testTicket := &cashew.Ticket{
		Type:             cashew.TicketTypeProxyGranting,
		LastReferencedAt: &testTime,
	}
	testee := New()

	err := testee.Validate(testTicket)
	if err != nil {
		t.Errorf("Validate() returned unexpected error %+v", err)
	}

	testTime = now.Add(time.Duration(-1*foundation.App().GrantingDefaultExpire-1) * time.Second)
	testTicketForError := &cashew.Ticket{
		Type:             cashew.TicketTypeProxyGranting,
		LastReferencedAt: &testTime,
	}

	err = testee.Validate(testTicketForError)
	if err != errors.ErrTicketHasBeenExpired {
		t.Errorf("Validate() returned unexpected error %+v", err)
	}

	// Confirm for nil LastReferencedAt
	testTicketForError.LastReferencedAt = nil
	err = testee.Validate(testTicketForError)
	if err != nil {
		t.Errorf("Validate() returned unexpected error %+v", err)
	}
}

func TestValidateTGT(t *testing.T) {
	loadSetting()
	now := time.Now()
	c := gomock.NewController(t)
	tm := timer.NewMockTimeWrapper(c)

	tm.EXPECT().Now().AnyTimes().Return(now)

	globaltimer.Local = tm

	testTime := now.Add(time.Duration(-1*foundation.App().GrantingDefaultExpire) * time.Second)

	testTicket := &cashew.Ticket{
		Type:             cashew.TicketTypeTicketGranting,
		LastReferencedAt: &testTime,
	}
	testee := New()

	err := testee.Validate(testTicket)
	if err != nil {
		t.Errorf("Validate() returned unexpected error %+v", err)
	}

	testTime = now.Add(time.Duration(-1*foundation.App().GrantingDefaultExpire-1) * time.Second)
	testTicketForError := &cashew.Ticket{
		Type:             cashew.TicketTypeTicketGranting,
		LastReferencedAt: &testTime,
	}

	err = testee.Validate(testTicketForError)
	if err != errors.ErrTicketHasBeenExpired {
		t.Errorf("Validate() returned unexpected error %+v", err)
	}

	// Confirm for nil LastReferencedAt
	testTicketForError.LastReferencedAt = nil
	err = testee.Validate(testTicketForError)
	if err != nil {
		t.Errorf("Validate() returned unexpected error %+v", err)
	}
}

func TestValidateOther(t *testing.T) {
	loadSetting()
	now := time.Now()
	c := gomock.NewController(t)
	tm := timer.NewMockTimeWrapper(c)

	tm.EXPECT().Now().AnyTimes().Return(now)

	globaltimer.Local = tm

	testTime := now.Add(time.Duration(-1*foundation.App().GrantingHardTimeout) * time.Second)

	testTicket := &cashew.Ticket{
		Type:             cashew.TicketTypeService,
		LastReferencedAt: &testTime,
		CreatedAt:        testTime,
	}
	testee := New()

	err := testee.Validate(testTicket)
	if err != nil {
		t.Errorf("Validate() returned unexpected error %+v", err)
	}

	testTime = now.Add(time.Duration(-1*foundation.App().GrantingHardTimeout-1) * time.Second)
	testTicketForError := &cashew.Ticket{
		Type:             cashew.TicketTypeProxy,
		LastReferencedAt: &testTime,
		CreatedAt:        testTime,
	}

	err = testee.Validate(testTicketForError)
	if err != errors.ErrHardTimeoutTicket {
		t.Errorf("Validate() returned unexpected error %+v", err)
	}
}

func TestValidateZeroConf(t *testing.T) {
	loadZeroSetting()
	now := time.Now()
	c := gomock.NewController(t)
	tm := timer.NewMockTimeWrapper(c)

	tm.EXPECT().Now().AnyTimes().Return(now)

	globaltimer.Local = tm

	testTime := now.Add(-1000000 * time.Second)

	testTicket := &cashew.Ticket{
		Type:      cashew.TicketTypeProxyGranting,
		CreatedAt: testTime,
	}
	testee := New()

	err := testee.Validate(testTicket)
	if err != nil {
		t.Errorf("Validate() returned unexpected error %+v", err)
	}

	testTicketForError := &cashew.Ticket{
		Type:      cashew.TicketTypeProxy,
		CreatedAt: testTime,
	}

	err = testee.Validate(testTicketForError)
	if err != nil {
		t.Errorf("Validate() returned unexpected error %+v", err)
	}
}

func loadSetting() {
	dir, _ := os.Getwd()
	path := filepath.Join(dir, "testconfig.yml")
	foundation.PrepareApp(path)
}

func loadZeroSetting() {
	dir, _ := os.Getwd()
	path := filepath.Join(dir, "zeroconfig.yml")
	foundation.PrepareApp(path)
}
