package channel

// SignalToInterference exports signalToInterference for testing.
func SignalToInterference(signal int) float64 {
	return signalToInterference(signal)
}

// GetChannelsForBand exports getChannelsForBand for testing.
func GetChannelsForBand(band Band) []int {
	return getChannelsForBand(band)
}

// GroupNetworksByChannel exports groupNetworksByChannel for testing.
func GroupNetworksByChannel(networks []NetworkInfo, band Band) map[int][]NetworkInfo {
	return groupNetworksByChannel(networks, band)
}

// BuildChannelInfo exports buildChannelInfo for testing.
func BuildChannelInfo(
	channel int,
	band Band,
	channelNetworks map[int][]NetworkInfo,
	allNetworks []NetworkInfo,
) ChannelInfo {
	return buildChannelInfo(channel, band, channelNetworks, allNetworks)
}

// FindBestChannel exports findBestChannel for testing.
func FindBestChannel(channels []ChannelInfo) int {
	return findBestChannel(channels)
}

// IsValid5GHzChannel exports isValid5GHzChannel for testing.
func IsValid5GHzChannel(channel int) bool {
	return isValid5GHzChannel(channel)
}

// Abs exports abs for testing.
func Abs(x int) int {
	return abs(x)
}
