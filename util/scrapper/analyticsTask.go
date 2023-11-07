package scrapper

import (
	"time"
)

type Channel struct {
	Name string
}

func RunAnalyticsTask(chRepo ChannelAuthorizerRepository, s *Scrapper) {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ticker.C:
			RunAnalytics(chRepo, *s)
		}
	}
}

func RunAnalytics(chRepo ChannelAuthorizerRepository, s Scrapper) {
	var channelList []*Channel
	channelList = chRepo.GetAllChannels()
	for _, channel := range channelList {
		s.CollectAvgViews(channel.Name)
	}

}
