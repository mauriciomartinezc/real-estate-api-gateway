package loadbalancer

var instances = []string{"http://instance1", "http://instance2"}
var currentIndex int

func GetNextInstance() string {
	instance := instances[currentIndex]
	currentIndex = (currentIndex + 1) % len(instances)
	return instance
}
