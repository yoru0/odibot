package dc

import "time"

func HandleShutdownCommand(id string) {
	if id == ownerID {
		go func() {
			time.Sleep(1 * time.Second)
			quit <- struct{}{}
		}()
	}
}
