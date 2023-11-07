package scrapper

import "fmt"

// DefaultRepo This default implementation of ChannelAuthorizerRepository should only be used for testing purposes
type DefaultRepo struct{}

func (r *DefaultRepo) SaveCode(c Code) error {
	fmt.Println(c.SubmitCode)
	return nil
}
func (r *DefaultRepo) SaveChannelAvgViews(ch string, v float64) error {
	fmt.Println(ch, v)
	return nil
}

func (r *DefaultRepo) GetAllChannels() []*Channel {
	var chArray []*Channel
	chArray[0] = &Channel{Name: "@smm_auto_test"}
	return chArray
}
