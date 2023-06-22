package union

type Union[T any, U any] struct {
	left  T
	right U
	set   int
}

func (union Union[T, U]) Unwrap() (*T, *U) {
	if union.IsLeft() {
		return &union.left, nil
	} else if union.IsRight() {
		return nil, &union.right
	} else {
		return nil, nil
	}
}

func (union Union[T, U]) ExpectLeft() T {
	return union.left
}

func (union Union[T, U]) ExpectRight() U {
	return union.right
}

func (union Union[T, U]) IsLeft() bool {
	return union.set == 1
}

func (union Union[T, U]) IsRight() bool {
	return union.set == 2
}

func (union Union[T, U]) Left(val T) Union[T, U] {
	union.set = 1
	union.left = val
	return union
}

func (union Union[T, U]) Right(val U) Union[T, U] {
	union.set = 2
	union.right = val
	return union
}
