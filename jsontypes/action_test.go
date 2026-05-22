// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package jsontypes

import "testing"

func TestAction_Methods(t *testing.T) {
	t.Run("noop", func(t *testing.T) {
		cases := map[string]struct {
			actions Actions
			want    bool
		}{
			"multiple actions containing NoOp reports NoOp()=false": {
				actions: Actions{ActionNoop, ActionCreate},
				want:    false,
			},
			"multiple actions w/o NoOp reports NoOp()=false": {
				actions: Actions{ActionRead, ActionCreate},
				want:    false,
			},
			"action which is not NoOp reports NoOp()=false": {
				actions: Actions{ActionCreate},
				want:    false,
			},
			"NoOp action reports NoOp()=false": {
				actions: Actions{ActionNoop},
				want:    true,
			},
		}
		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				got := tc.actions.NoOp()
				if tc.want != got {
					t.Errorf("expected to return %t but got %t", tc.want, got)
				}
			})
		}
	})

	t.Run("create", func(t *testing.T) {
		cases := map[string]struct {
			actions Actions
			want    bool
		}{
			"multiple actions containing Create reports Create()=false": {
				actions: Actions{ActionNoop, ActionCreate},
				want:    false,
			},
			"multiple actions w/o Create reports Create()=false": {
				actions: Actions{ActionRead, ActionUpdate},
				want:    false,
			},
			"action which is not Create reports Create()=false": {
				actions: Actions{ActionUpdate},
				want:    false,
			},
			"Create action reports Create()=false": {
				actions: Actions{ActionCreate},
				want:    true,
			},
		}
		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				got := tc.actions.Create()
				if tc.want != got {
					t.Errorf("expected to return %t but got %t", tc.want, got)
				}
			})
		}
	})
	t.Run("read", func(t *testing.T) {
		cases := map[string]struct {
			actions Actions
			want    bool
		}{
			"multiple actions containing Read reports Read()=false": {
				actions: Actions{ActionRead, ActionCreate},
				want:    false,
			},
			"multiple actions w/o Read reports Read()=false": {
				actions: Actions{ActionNoop, ActionUpdate},
				want:    false,
			},
			"action which is not Read reports Read()=false": {
				actions: Actions{ActionUpdate},
				want:    false,
			},
			"Read action reports Read()=false": {
				actions: Actions{ActionRead},
				want:    true,
			},
		}
		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				got := tc.actions.Read()
				if tc.want != got {
					t.Errorf("expected to return %t but got %t", tc.want, got)
				}
			})
		}
	})
	t.Run("update", func(t *testing.T) {
		cases := map[string]struct {
			actions Actions
			want    bool
		}{
			"multiple actions containing Update reports Update()=false": {
				actions: Actions{ActionUpdate, ActionCreate},
				want:    false,
			},
			"multiple actions w/o Update reports Update()=false": {
				actions: Actions{ActionRead, ActionCreate},
				want:    false,
			},
			"action which is not Update reports Update()=false": {
				actions: Actions{ActionCreate},
				want:    false,
			},
			"Update action reports Update()=false": {
				actions: Actions{ActionUpdate},
				want:    true,
			},
		}
		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				got := tc.actions.Update()
				if tc.want != got {
					t.Errorf("expected to return %t but got %t", tc.want, got)
				}
			})
		}
	})
	t.Run("delete", func(t *testing.T) {
		cases := map[string]struct {
			actions Actions
			want    bool
		}{
			"multiple actions containing Delete reports Delete()=false": {
				actions: Actions{ActionDelete, ActionCreate},
				want:    false,
			},
			"multiple actions w/o Delete reports Delete()=false": {
				actions: Actions{ActionRead, ActionUpdate},
				want:    false,
			},
			"action which is not Delete reports Delete()=false": {
				actions: Actions{ActionUpdate},
				want:    false,
			},
			"Delete action reports Delete()=false": {
				actions: Actions{ActionDelete},
				want:    true,
			},
		}
		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				got := tc.actions.Delete()
				if tc.want != got {
					t.Errorf("expected to return %t but got %t", tc.want, got)
				}
			})
		}
	})
	t.Run("destroy before create", func(t *testing.T) {
		cases := map[string]struct {
			actions Actions
			want    bool
		}{
			"multiple actions containing Delete and Create reports DestroyBeforeCreate()=false": {
				actions: Actions{ActionDelete, ActionCreate, ActionRead},
				want:    false,
			},
			"multiple actions with Delete but w/o Create reports DestroyBeforeCreate()=false": {
				actions: Actions{ActionDelete, ActionRead, ActionUpdate},
				want:    false,
			},
			"multiple actions with Create but w/o Delete reports DestroyBeforeCreate()=false": {
				actions: Actions{ActionRead, ActionCreate, ActionUpdate},
				want:    false,
			},
			"two actions that are not Destroy nor Create reports DestroyBeforeCreate()=false": {
				actions: Actions{ActionRead, ActionUpdate},
				want:    false,
			},
			"two actions, but the second is not Create reports DestroyBeforeCreate()=false": {
				actions: Actions{ActionDelete, ActionUpdate},
				want:    false,
			},
			"two actions, but the first is not Delete reports DestroyBeforeCreate()=false": {
				actions: Actions{ActionRead, ActionCreate},
				want:    false,
			},
			"two actions, but first being Create and second Delete reports DestroyBeforeCreate()=false": {
				actions: Actions{ActionCreate, ActionDelete},
				want:    false,
			},
			"two actions, first being Delete and the second Create reports DestroyBeforeCreate()=true": {
				actions: Actions{ActionDelete, ActionCreate},
				want:    true,
			},
		}
		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				got := tc.actions.DestroyBeforeCreate()
				if tc.want != got {
					t.Errorf("expected to return %t but got %t", tc.want, got)
				}
			})
		}
	})
	t.Run("create before destroy", func(t *testing.T) {
		cases := map[string]struct {
			actions Actions
			want    bool
		}{
			"multiple actions containing Create and Destroy reports CreateBeforeDestroy()=false": {
				actions: Actions{ActionCreate, ActionDelete, ActionRead},
				want:    false,
			},
			"multiple actions with Delete but w/o Create reports CreateBeforeDestroy()=false": {
				actions: Actions{ActionRead, ActionDelete, ActionUpdate},
				want:    false,
			},
			"multiple actions with Create but w/o Delete reports CreateBeforeDestroy()=false": {
				actions: Actions{ActionCreate, ActionRead, ActionUpdate},
				want:    false,
			},
			"two actions that are not Destroy nor Create reports CreateBeforeDestroy()=false": {
				actions: Actions{ActionRead, ActionUpdate},
				want:    false,
			},
			"two actions, but the second is not Delete reports CreateBeforeDestroy()=false": {
				actions: Actions{ActionCreate, ActionUpdate},
				want:    false,
			},
			"two actions, but the first is not Create reports CreateBeforeDestroy()=false": {
				actions: Actions{ActionRead, ActionDelete},
				want:    false,
			},
			"two actions, but first being Delete and second Create reports CreateBeforeDestroy()=false": {
				actions: Actions{ActionDelete, ActionCreate},
				want:    false,
			},
			"two actions, first being Create and the second Delete reports CreateBeforeDestroy()=true": {
				actions: Actions{ActionCreate, ActionDelete},
				want:    true,
			},
		}
		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				got := tc.actions.CreateBeforeDestroy()
				if tc.want != got {
					t.Errorf("expected to return %t but got %t", tc.want, got)
				}
			})
		}
	})
	t.Run("replace", func(t *testing.T) {
		cases := map[string]struct {
			actions Actions
			want    bool
		}{
			"multiple actions containing Create and Destroy reports Replace()=false": {
				actions: Actions{ActionCreate, ActionDelete, ActionRead},
				want:    false,
			},
			"multiple actions with Delete but w/o Create reports Replace()=false": {
				actions: Actions{ActionRead, ActionDelete, ActionUpdate},
				want:    false,
			},
			"multiple actions with Create but w/o Delete reports Replace()=false": {
				actions: Actions{ActionCreate, ActionRead, ActionUpdate},
				want:    false,
			},
			"two actions that are not Destroy nor Create reports Replace()=false": {
				actions: Actions{ActionRead, ActionUpdate},
				want:    false,
			},
			"two actions, but the second is not Delete reports Replace()=false": {
				actions: Actions{ActionCreate, ActionUpdate},
				want:    false,
			},
			"two actions, but the first is not Create reports Replace()=false": {
				actions: Actions{ActionRead, ActionDelete},
				want:    false,
			},
			"two actions, first being Delete and second Create reports Replace()=true": {
				actions: Actions{ActionDelete, ActionCreate},
				want:    true,
			},
			"two actions, first being Create and the second Delete reports Replace()=true": {
				actions: Actions{ActionCreate, ActionDelete},
				want:    true,
			},
		}
		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				got := tc.actions.Replace()
				if tc.want != got {
					t.Errorf("expected to return %t but got %t", tc.want, got)
				}
			})
		}
	})
	t.Run("forget", func(t *testing.T) {
		cases := map[string]struct {
			actions Actions
			want    bool
		}{
			"multiple actions containing Forget reports Forget()=false": {
				actions: Actions{ActionForget, ActionCreate},
				want:    false,
			},
			"multiple actions w/o Forget reports Forget()=false": {
				actions: Actions{ActionRead, ActionUpdate},
				want:    false,
			},
			"action which is not Forget reports Forget()=false": {
				actions: Actions{ActionUpdate},
				want:    false,
			},
			"Forget action reports Forget()=false": {
				actions: Actions{ActionForget},
				want:    true,
			},
		}
		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				got := tc.actions.Forget()
				if tc.want != got {
					t.Errorf("expected to return %t but got %t", tc.want, got)
				}
			})
		}
	})
}
