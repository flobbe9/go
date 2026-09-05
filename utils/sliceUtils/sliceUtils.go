package sliceUtils

// Classic "map" function.
//
// [s] the slice to map
//
// [mapperCallback] the callback. Accepts the element of iteration and the index. Returns an element of type [N]
//
// <S> the type of slice
//
// <N> the new type of slice beeing returned
//
// [return] new slice containing an equivalent element of type [N] for each element from [s]
//
// [panic] if an arg is [nil]
func Map[S, N any](s []S, mapperCallback func(element S, index int) N) []N {
	if s == nil || mapperCallback == nil {
		panic("Nil arg");
	}

	newSlice := []N{};
	for i, element := range s {
		newSlice = append(newSlice, mapperCallback(element, i));
	} 

	return newSlice;
}

// Classic "filter" function.
//
// [s] the slice to filter
//
// [predicate] the callback used on every element in [s]. Accepts the element of iteration and the index and returns [true]
// if the element should not be filtered out, else [false].
//
// <S> the type of slice
//
// [return] new slice containing a subset of [s]
//
// [panic] if an arg is [nil]
func Filter[S any](s []S, predicate func(element S, index int) bool) []S {
	if s == nil || predicate == nil {
		panic("Nil arg");
	}

	filtered := []S{};

	for i, element := range s {
		if predicate(element, i) {
			filtered = append(filtered, element);
		}
	}

	return filtered;
}