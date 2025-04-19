package reward

import "context"

type Repository interface {
	SaveReward(ctx context.Context, reward *Reward) error
}
