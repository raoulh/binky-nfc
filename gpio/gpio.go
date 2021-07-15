package gpio

import (
	"os"
	"io/ioutil"
)

type GPIO struct{}

func (r GPIO) Pin(name string) GPIO_Pin {
	pin := GPIO_Pin{name}
	filename := pin.Filename()
	if _, err := os.Stat(filename); os.IsNotExist(err) {
			// export gpio pin
			ioutil.WriteFile("/sys/class/gpio/export", []byte(pin.Name), 0666)
	}
	return pin
}

type GPIO_Pin struct {
	Name string
}

func (r GPIO_Pin) Filename() string {
	return "/sys/class/gpio/gpio" + r.Name
}

func (r GPIO_Pin) write(where, what string) GPIO_Pin {
	filename := r.Filename() + "/" + where
	ioutil.WriteFile(filename, []byte(what), 0666)
	return r
}

func (r GPIO_Pin) Output() GPIO_Pin {
	return r.write("direction", "out")
}

func (r GPIO_Pin) High() GPIO_Pin {
	return r.write("value", "1")
}

func (r GPIO_Pin) Low() GPIO_Pin {
	return r.write("value", "0")
}
