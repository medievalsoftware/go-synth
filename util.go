package snd

import "bufio"

func smart(reader *bufio.Reader) (int, error) {
	if data, err := reader.ReadByte(); err != nil {
		return 0, err
	} else if (data & 0b1000_0000) == 0 {
		return (int(data) & 0xFF) - 64, nil
	} else {
		val := (int(data) & 0xFF) << 8

		if data, err := reader.ReadByte(); err != nil {
			return 0, err
		} else {
			val |= int(data) & 0xFF
			val -= 49152
			return val, nil
		}
	}
}

func usmart(reader *bufio.Reader) (int, error) {
	if data, err := reader.ReadByte(); err != nil {
		return 0, err
	} else if (data & 0b1000_0000) == 0 {
		return int(data) & 0xFF, nil
	} else {
		val := (int(data) & 0xFF) << 8

		if data, err := reader.ReadByte(); err != nil {
			return 0, err
		} else {
			val |= int(data) & 0xFF
			val -= 32768
			return val, nil
		}
	}
}
