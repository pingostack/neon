package rtc

import (
	"fmt"
	"strconv"
)

func (p Parameters) Integer(key string) (int64, error) {
	val, found := p[key]
	if !found {
		return 0, fmt.Errorf("key not found: %s", key)
	}

	return strconv.ParseInt(val, 10, 64)
}

func (p Parameters) String(key string) (string, error) {
	val, found := p[key]
	if !found {
		return "", fmt.Errorf("key not found: %s", key)
	}

	return val, nil
}

func (p Parameters) Bool(key string) (bool, error) {
	val, found := p[key]
	if !found {
		return false, fmt.Errorf("key not found: %s", key)
	}

	return strconv.ParseBool(val)
}

func (p Parameters) Float(key string) (float64, error) {
	val, found := p[key]
	if !found {
		return 0, fmt.Errorf("key not found: %s", key)
	}

	return strconv.ParseFloat(val, 64)
}

func (p Parameters) Uint8(key string) (uint8, error) {
	val, found := p[key]
	if !found {
		return 0, fmt.Errorf("key not found: %s", key)
	}

	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint8(i), nil
}

func (p Parameters) Set(key string, value interface{}) {
	p[key] = fmt.Sprintf("%v", value)
}

func (p Parameters) Delete(key string) {
	delete(p, key)
}

func (p Parameters) Has(key string) bool {
	_, found := p[key]
	return found
}

func (p Parameters) Clone() Parameters {
	clone := Parameters{}
	for k, v := range p {
		clone[k] = v
	}
	return clone
}

func (p Parameters) Merge(other Parameters) {
	for k, v := range other {
		p[k] = v
	}
}

func (p Parameters) MergeOverwrite(other Parameters) {
	for k, v := range other {
		if _, found := p[k]; !found {
			p[k] = v
		}
	}
}
