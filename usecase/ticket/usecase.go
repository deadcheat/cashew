package ticket

import (
	"database/sql"
	"net/http"
	"net/url"

	"github.com/deadcheat/cashew"
	"github.com/deadcheat/cashew/errors"
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

// Find find ticket by id
func (u *UseCase) Find(id string) (t *cashew.Ticket, err error) {
	if id == "" {
		return nil, errors.ErrNoTicketID
	}
	t, err = u.r.Find(id)
	if err == sql.ErrNoRows {
		err = errors.ErrTicketNotFound
	}
	return
}

// NewLogin create new LoginTicket
func (u *UseCase) NewLogin(r *http.Request) (t *cashew.Ticket, err error) {
	t = new(cashew.Ticket)
	t.Type = cashew.TicketTypeLogin
	t.ID = u.idr.Issue(t.Type)
	t.ClientHostName = u.chr.Ensure(r)
	if err = u.r.Create(t); err != nil {
		return nil, err
	}

	return
}

// NewProxyGranting create new proxy-granting-ticket
func (u *UseCase) NewProxyGranting(r *http.Request, callbackURL *url.URL, st *cashew.Ticket) (t *cashew.Ticket, err error) {
	// if callback-url is not set, skip this method
	if callbackURL == nil {
		return nil, nil
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

// NewService create new Service-Ticket
func (u *UseCase) NewService(r *http.Request, service *url.URL, tgt *cashew.Ticket, primary bool) (t *cashew.Ticket, err error) {
	if service == nil {
		return nil, errors.ErrNoServiceDetected
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

// NewProxy create new ProxyTicket
func (u *UseCase) NewProxy(r *http.Request, service string, grantedBy *cashew.Ticket) (t *cashew.Ticket, err error) {
	if service == "" {
		return nil, errors.ErrNoServiceDetected
	}
	t = new(cashew.Ticket)
	t.Type = cashew.TicketTypeProxy
	t.ID = u.idr.Issue(t.Type)
	t.ClientHostName = u.chr.Ensure(r)
	t.Service = service
	t.UserName = grantedBy.UserName
	t.GrantedBy = grantedBy
	if err = u.r.Create(t); err != nil {
		return nil, err
	}
	return
}

// NewGranting create new ServiceTicket
func (u *UseCase) NewGranting(r *http.Request, username string, extraAttributes interface{}) (t *cashew.Ticket, err error) {
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

// Consume will consume and update 'last-referenced-at'
func (u *UseCase) Consume(t *cashew.Ticket) error {
	return u.r.Consume(t)
}
