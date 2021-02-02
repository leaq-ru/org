package consumer

import "time"

func (c Consumer) pollSubIsValid() (err error) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.state.done:
			return
		case <-ticker.C:
			if c.state.sub.IsValid() {
				continue
			}

			err = c.Subscribe()
			if err != nil {
				c.logger.Error().Err(err).Send()
				return
			}
		}
	}
}
