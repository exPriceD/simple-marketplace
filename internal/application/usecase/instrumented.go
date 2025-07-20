package usecase

import (
	"context"
	"time"

	"github.com/exPriceD/simple-marketplace/internal/application/dto"
	apperror "github.com/exPriceD/simple-marketplace/internal/application/error"
	"github.com/exPriceD/simple-marketplace/internal/platform/log"
	"github.com/exPriceD/simple-marketplace/internal/platform/metrics"
)

func statusFromError(err error) string {
	if err == nil {
		return "ok"
	}
	if ae, ok := err.(*apperror.AppError); ok {
		return string(ae.Kind)
	}
	return "error"
}

type registerUserInstrumented struct {
	inner  RegisterUserUseCase
	logger log.Logger
	m      metrics.Metrics
	op     string
}

func InstrumentRegisterUser(inner RegisterUserUseCase, logger log.Logger, m metrics.Metrics) RegisterUserUseCase {
	return &registerUserInstrumented{inner: inner, logger: logger, m: m, op: "user.register"}
}

func (w *registerUserInstrumented) Execute(ctx context.Context, in RegisterUserInput) (dto.UserDTO, error) {
	start := time.Now()
	out, err := w.inner.Execute(ctx, in)
	el := time.Since(start).Seconds()
	st := statusFromError(err)
	w.m.IncUseCaseCalls(w.op, st)
	w.m.ObserveUseCaseDuration(w.op, el)
	w.logger.Info(ctx, "usecase", "op", w.op, "status", st, "duration_ms", int64(el*1000))
	return out, err
}

type loginUserInstrumented struct {
	inner  LoginUserUseCase
	logger log.Logger
	m      metrics.Metrics
	op     string
}

func InstrumentLoginUser(inner LoginUserUseCase, logger log.Logger, m metrics.Metrics) LoginUserUseCase {
	return &loginUserInstrumented{inner: inner, logger: logger, m: m, op: "user.login"}
}

func (w *loginUserInstrumented) Execute(ctx context.Context, in LoginUserInput) (dto.AuthDTO, error) {
	start := time.Now()
	out, err := w.inner.Execute(ctx, in)
	el := time.Since(start).Seconds()
	st := statusFromError(err)
	w.m.IncUseCaseCalls(w.op, st)
	w.m.ObserveUseCaseDuration(w.op, el)
	w.logger.Info(ctx, "usecase", "op", w.op, "status", st, "duration_ms", int64(el*1000))
	return out, err
}

type createListingInstrumented struct {
	inner  CreateListingUseCase
	logger log.Logger
	m      metrics.Metrics
	op     string
}

func InstrumentCreateListing(inner CreateListingUseCase, logger log.Logger, m metrics.Metrics) CreateListingUseCase {
	return &createListingInstrumented{inner: inner, logger: logger, m: m, op: "listing.create"}
}

func (w *createListingInstrumented) Execute(ctx context.Context, in CreateListingInput) (dto.ListingDTO, error) {
	start := time.Now()
	out, err := w.inner.Execute(ctx, in)
	el := time.Since(start).Seconds()
	st := statusFromError(err)
	w.m.IncUseCaseCalls(w.op, st)
	w.m.ObserveUseCaseDuration(w.op, el)
	w.logger.Info(ctx, "usecase", "op", w.op, "status", st, "duration_ms", int64(el*1000))
	return out, err
}

type listListingsInstrumented struct {
	inner  ListListingsUseCase
	logger log.Logger
	m      metrics.Metrics
	op     string
}

func InstrumentListListings(inner ListListingsUseCase, logger log.Logger, m metrics.Metrics) ListListingsUseCase {
	return &listListingsInstrumented{inner: inner, logger: logger, m: m, op: "listing.list"}
}

func (w *listListingsInstrumented) Execute(ctx context.Context, in ListListingsInput) (dto.ListingFeedDTO, error) {
	start := time.Now()
	out, err := w.inner.Execute(ctx, in)
	el := time.Since(start).Seconds()
	st := statusFromError(err)
	w.m.IncUseCaseCalls(w.op, st)
	w.m.ObserveUseCaseDuration(w.op, el)
	w.logger.Info(ctx, "usecase", "op", w.op, "status", st, "duration_ms", int64(el*1000))
	return out, err
}
