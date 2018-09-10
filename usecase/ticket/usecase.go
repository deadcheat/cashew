package ticket

import (
	"net/http"
	"net/url"

	"github.com/deadcheat/cashew"
	"github.com/deadcheat/cashew/values/errs"
)

// UseCase implemented cashew.LoginUseCase
type UseCase struct {
	r   cashew.TicketRepository
	idr cashew.IDRepository
	chr cashew.ClientHostNameRepository
	pcr cashew.ProxyCallBackRepository
}

// New return new usecase
func New(r cashew.TicketRepository, idr cashew.IDRepository, chr cashew.ClientHostNameRepository, pcr cashew.ProxyCallBackRepository) cashew.TicketUseCase {
	return &UseCase{r, idr, chr, pcr}
}

// FindTicket find ticket by id
func (u *UseCase) FindTicket(id string) (*cashew.Ticket, error) {
	if id == "" {
		return nil, errs.ErrNoTicketID
	}
	return u.r.Find(id)
}

// LoginTicket create new LoginTicket
func (u *UseCase) LoginTicket(r *http.Request) (t *cashew.Ticket, err error) {
	t = new(cashew.Ticket)
	t.Type = cashew.TicketTypeLogin
	t.ID = u.idr.Issue(t.Type)
	t.ClientHostName = u.chr.Ensure(r)
	if err = u.r.Create(t); err != nil {
		return nil, err
	}

	return
}

// ProxyGrantingTicket create new ProxyGrantingTicket
func (u *UseCase) ProxyGrantingTicket(r *http.Request, callbackURL *url.URL, st *cashew.Ticket) (t *cashew.Ticket, err error) {
	if callbackURL == nil {
		return nil, errs.ErrProxyCallBackURLMissing
	}
	t = new(cashew.Ticket)
	t.Type = cashew.TicketTypeProxyGranting
	t.ID = u.idr.Issue(t.Type)
	t.IOU = u.idr.Issue(cashew.TicketTypeProxyGrantingIOU)
	t.ClientHostName = u.chr.Ensure(r)
	t.GrantedBy = st

	if err = u.pcr.Dial(callbackURL, t.ID, t.IOU); err != nil {
		return nil, err
	}

	if err = u.r.Create(t); err != nil {
		return
	}
	return
}

// ServiceTicket create new ServiceTicket
func (u *UseCase) ServiceTicket(r *http.Request, service *url.URL, tgt *cashew.Ticket, primary bool) (t *cashew.Ticket, err error) {
	if service == nil {
		return nil, errs.ErrNoServiceDetected
	}
	t = new(cashew.Ticket)
	t.Type = cashew.TicketTypeService
	t.ID = u.idr.Issue(t.Type)
	t.ClientHostName = u.chr.Ensure(r)
	t.Service = service.String()
	t.UserName = tgt.UserName
	t.GrantedBy = tgt
	t.Primary = primary
	if err = u.r.Create(t); err != nil {
		return nil, err
	}

	return
}

// TicketGrantingTicket create new ServiceTicket
func (u *UseCase) TicketGrantingTicket(r *http.Request, username string, extraAttributes interface{}) (t *cashew.Ticket, err error) {
	t = new(cashew.Ticket)
	t.Type = cashew.TicketTypeTicketGranting
	t.ID = u.idr.Issue(t.Type)
	t.ClientHostName = u.chr.Ensure(r)
	t.UserName = username
	t.ExtraAttributes = extraAttributes
	if err = u.r.Create(t); err != nil {
		return nil, err
	}

	return
}
