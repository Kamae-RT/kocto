package kocto

import "time"

// PageStream streams data or errors out of `worker` into channel.
//
// It assumes `worker` is able to calculate if it has more pages,
// and pages are numeric
func PageStream[T any](errs chan error, worker func(int) ([]T, bool, error)) chan T {
	stream := make(chan T)

	go func() {
		hasMore := true
		var page int = 1

		for hasMore {
			data, hm, err := worker(page)
			if err != nil {
				errs <- err
				break
			}

			for _, datum := range data {
				stream <- datum
			}

			hasMore = hm
			page = page + 1
		}

		close(stream)
	}()

	return stream
}

// LinkStream produces a stream of T from `worker` into channel.
//
// It assumes `worker` is able to provide the next page link and accepts an empty page.
func LinkStream[T any](errs chan error, worker func(string) ([]T, string, error)) chan T {
	data := make(chan T)

	go func() {
		hasMore := true
		nextPage := ""

		for hasMore {
			ds, next, err := worker(nextPage)
			if err != nil {
				errs <- err
				break
			}

			for _, d := range ds {
				data <- d
			}

			hasMore = next != ""
			nextPage = next
		}

		close(data)
	}()

	return data
}

type Range struct {
	Start time.Time
	End   time.Time
}

const Month = time.Hour * 24 * 30

// RangeLinkStream produces a stream of T from worker
func RangeLinkStream[T any](errs chan error, r Range, pf func(Range, string) ([]T, string, error)) chan T {
	data := make(chan T)

	go func() {
		max := r.Start
		currentEnd := r.End
		streamErrs := make(chan error)

	dateLoop:
		for max.Before(currentEnd) {
			var currentStart = currentEnd.Add(-Month)
			if currentStart.Before(max) {
				currentStart = max
			}

			month := Range{
				Start: currentStart,
				End:   currentEnd,
			}

			monthData := LinkStreamWithRange(streamErrs, month, pf)

			for monthData != nil {
				select {
				case err := <-streamErrs:
					errs <- err
					break dateLoop

				case d, ok := <-monthData:
					if !ok {
						monthData = nil
						continue
					}

					data <- d
				}
			}

			currentEnd = currentStart
		}

		close(data)
	}()

	return data
}

// LinkStreamWithRange produces a stream of T from a worker that accepts a range 
//
// It assumes `worker` is able to provide the next page link and accepts an empty page.
func LinkStreamWithRange[T any](errs chan error, r Range, worker func(Range, string) ([]T, string, error)) chan T {
	stream := make(chan T)

	go func() {
		hasMore := true
		nextPage := ""

		for hasMore {
			ds, next, err := worker(r, nextPage)
			if err != nil {
				errs <- err
				break
			}

			for _, d := range ds {
				stream <- d
			}

			hasMore = next != ""
			nextPage = next
		}

		close(stream)
	}()

	return stream
}
