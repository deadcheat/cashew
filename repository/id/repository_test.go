package id

import (
	"testing"

	"github.com/deadcheat/cashew"
	"github.com/deadcheat/cashew/values/consts"
)

func TestIssue(t *testing.T) {
	r := New()
	if tid := r.Issue(cashew.TicketTypeLogin); tid == "" {
		t.Error("this should return some string")
	}
}

func TestPrefixPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
		t.Error("this should catch error")
	}()
	dummy := cashew.TicketType(100)

	prefix(dummy)
}

func TestPrefix(t *testing.T) {
	type args struct {
		t cashew.TicketType
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: consts.TicketTypeStrLogin,
			args: args{cashew.TicketTypeLogin},
			want: consts.PrefixLoginTicket,
		},
		{
			name: consts.TicketTypeStrService,
			args: args{cashew.TicketTypeService},
			want: consts.PrefixServiceTicket,
		},
		{
			name: consts.TicketTypeStrProxy,
			args: args{cashew.TicketTypeProxy},
			want: consts.PrefixProxyTicket,
		},
		{
			name: consts.TicketTypeStrTicketGranting,
			args: args{cashew.TicketTypeTicketGranting},
			want: consts.PrefixTicketGrantingCookie,
		},
		{
			name: consts.TicketTypeStrProxyGranting,
			args: args{cashew.TicketTypeProxyGranting},
			want: consts.PrefixProxyGrantingTicket,
		},
		{
			name: consts.TicketTypeStrProxyGrantingIOU,
			args: args{cashew.TicketTypeProxyGrantingIOU},
			want: consts.PrefixProxyGrantingTicketIOU,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := prefix(tt.args.t); got != tt.want {
				t.Errorf("prefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
