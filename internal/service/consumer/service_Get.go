package consumer

import (
	"context"
	"errors"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore/consumers"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
)

func (s *service) Get(ctx context.Context, input GetInput) (output GetOutput, err error) {
	consumerOutput, err := s.consumerRepository.Get(ctx, consumers.GetInput{
		ID:     input.ConsumerID,
		UserID: input.UserID,
	})
	if err != nil {
		if errors.Is(err, datastore.ErrRecordNotFound) {
			err = ErrConsumerNotFound
		}
		return output, tracer.Error(err)
	}

	output = GetOutput{
		ID:        consumerOutput.ID,
		UserID:    consumerOutput.UserID,
		FullName:  consumerOutput.FullName,
		LegalName: consumerOutput.LegalName,
	}
	return
}
