package generics

func Map[U []F, T ~[]E, E any, F any](items T, callback func(value E, index int) F) U {
	mapped := make(U, len(items), cap(items))
	for index, item := range items {
		mapped[index] = callback(item, index)
	}
	return mapped
}

func Contains[E comparable, T ~[]E](items T, value E) bool {
	for _, v := range items {
		if v == value {
			return true
		}
	}
	return false
}

func Filter[T ~[]E, E any](items T, callback func(value E, index int) bool) T {
	var filtered = make(T, 0)
	for index, item := range items {
		if callback(item, index) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func Remove[T ~[]E, E comparable](items T, item E) T {
	var rit = make(T, 0)
	for _, it := range items {
		if item != it {
			rit = append(rit, it)
		}
	}
	return rit
}

func RemoveWhen[T ~[]E, E any](items T, item E, when func(item, value E, index int) bool) T {
	var rit = make(T, 0)
	for i, it := range items {
		if when(it, item, i) {
			rit = append(rit, it)
		}
	}
	return rit
}

func Match[T ~[]E, E any](items T, when func(item E, index int) bool) bool {
	for i, v := range items {
		if when(v, i) {
			return true
		}
	}
	return false
}

func Find[T ~[]E, E any](items T, when func(item E, index int) bool) int {
	for i, i2 := range items {
		if when(i2, i) {
			return i
		}
	}
	return -1
}
