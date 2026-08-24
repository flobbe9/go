package slicesUtils;


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
func SlicesMap[S any, N any](s []S, mapperCallback func(element S, index int) N) []N {
	if s == nil || mapperCallback == nil {
		panic("Nil arg");
	}

	newSlice := []N{};
	for i, element := range s {
		newSlice = append(newSlice, mapperCallback(element, i));
	} 

	return newSlice;
}