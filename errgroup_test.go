package errgroup

import (
	"context"
	"errors"
	"testing"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name string
		errs []error
		want string
	}{
		{
			name: "return 'err1;err2' when the errs are 'err1' and 'err2' and nil",
			errs: []error{
				errors.New("err1"),
				nil,
				errors.New("err2"),
			},
			want: "err1;err2",
		},
		{
			name: "return 'err1;err2' when the errs is nil",
			errs: []error{
				nil,
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Error{
				errs: tt.errs,
			}.Error()
			if got != tt.want {
				t.Errorf("want %s, but got: %s", tt.want, got)
			}
		})
	}
}

func TestError_Errors(t *testing.T) {
	type test struct {
		name string
		errs []error
		want []error
	}
	tests := []test{
		func() test {
			errs := []error{
				errors.New("err1"),
				errors.New("err2"),
			}
			return test{
				name: "return 'err1' and 'err2' when the errs are 'err1' and 'err2'",
				errs: errs,
				want: errs,
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Error{
				errs: tt.errs,
			}.Errors()
			if len(got) != len(tt.want) {
				t.Errorf("len want %d, but got %d", len(tt.want), len(got))
			}
			for i := range tt.want {
				if tt.want[i].Error() != got[i].Error() {
					t.Errorf("want %v, but got %v", tt.want[i], got[i])
				}
			}
		})
	}
}

func TestGroup_Wait(t *testing.T) {
	type test struct {
		name       string
		want       error
		init       func() *Group
		beforeFunc func(*Group)
		checkFunc  func() error
	}

	tests := []test{
		func() test {
			errs := []error{
				errors.New("err1"),
				errors.New("err2"),
				errors.New("err3"),
			}
			return test{
				name: "",
				init: func() *Group {
					return new(Group)
				},
				want: Error{
					errs: errs,
				},
				beforeFunc: func(g *Group) {
					for _, err := range append(errs, nil, nil) {
						err := err
						g.Go(func() error {
							return err
						})
					}
				},
			}
		}(),

		func() test {
			errs := []error{
				errors.New("err1"),
				errors.New("err2"),
				errors.New("err3"),
			}
			var ctx context.Context

			return test{
				name: "",
				init: func() (eg *Group) {
					eg, ctx = WithContext(context.Background())
					return
				},
				want: Error{
					errs: errs,
				},
				beforeFunc: func(g *Group) {
					for _, err := range append(errs, nil, nil) {
						err := err
						g.Go(func() error {
							return err
						})
					}
				},
				checkFunc: func() error {
					if ctx.Err() != context.Canceled {
						return errors.New("eg.ctx.Error() is not canceled")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.init == nil {
				t.Error("initialize function not set")
			}

			eg := tt.init()

			if tt.beforeFunc != nil {
				tt.beforeFunc(eg)
			}

			got := eg.Wait()

			if tt.checkFunc != nil {
				if err := tt.checkFunc(); err != nil {
					t.Errorf("checkFunc error: %v", err)
				}
			}

			if (tt.want == nil && got != nil) || (tt.want != nil && got == nil) {
				t.Errorf("err want %v, but got %v", tt.want, got)
			} else {
				got, ok := got.(Error)
				if !ok {
					t.Error("err type is not Error")
				}

				want := tt.want.(Error)

				if len(got.Errors()) != len(want.Errors()) {
					t.Errorf("len(want.Errors()) %d, but len(got.Errors()) %d", len(got.Errors()), len(got.Errors()))
				}

				for _, werr := range want.Errors() {
					var ok bool
					for _, gerr := range got.Errors() {
						if werr == gerr {
							ok = true
							break
						}
					}
					if !ok {
						t.Errorf("want %v, but got %v", want, got)
					}
				}
			}
		})
	}
}
