package consumption_form_page

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/stretchr/testify/assert"
	"ktea/kadmin"
	"ktea/kontext"
	"ktea/tests/keys"
	"ktea/ui"
	"ktea/ui/pages/nav"
	"testing"
)

func TestConsumeForm_Navigation(t *testing.T) {

	t.Run("esc goes back to topic list page", func(t *testing.T) {
		m := New(&kadmin.ListedTopic{
			Name:           "topic1",
			Replicas:       1,
			PartitionCount: 10,
		}, ui.NewTestKontext())
		// make sure form has been initialized
		m.View(ui.NewTestKontext(), ui.TestRenderer)

		cmd := m.Update(keys.Key(tea.KeyEsc))

		assert.IsType(t, nav.LoadTopicsPageMsg{}, cmd())
	})

	t.Run("renders all available partitions when there is height enough", func(t *testing.T) {
		m := New(&kadmin.ListedTopic{
			Name:           "topic1",
			Replicas:       1,
			PartitionCount: 10,
		}, ui.NewTestKontext())

		// make sure form has been initialized
		m.View(ui.TestKontext, ui.TestRenderer)

		render := m.View(ui.NewTestKontext(), ui.TestRenderer)

		assert.Contains(t, render, "> • 0")
		for i := 1; i < 10; i++ {
			assert.Regexp(t, fmt.Sprintf("• %d", i), render)
		}
		assert.NotContains(t, render, "• 10")
	})

	t.Run("renders subset of partitions when there is not enough height", func(t *testing.T) {
		ktx := &kontext.ProgramKtx{
			Config:          nil,
			WindowWidth:     100,
			WindowHeight:    20,
			AvailableHeight: 20,
		}
		m := New(&kadmin.ListedTopic{
			Name:           "topic1",
			Replicas:       1,
			PartitionCount: 100,
		}, ktx)
		// make sure form has been initialized
		m.View(ktx, ui.TestRenderer)

		render := m.View(ktx, ui.TestRenderer)

		assert.Contains(t, render, `> • 0`)
		assert.Contains(t, render, `• 1`)
		assert.Contains(t, render, `• 2`)
		assert.Contains(t, render, `• 3`)
		assert.Contains(t, render, `• 4`)
		assert.NotContains(t, render, "• 5")
	})

	t.Run("load form based on previous ReadDetails", func(t *testing.T) {
		m := NewWithDetails(&kadmin.ReadDetails{
			TopicName:       "topic1",
			PartitionToRead: []int{3, 6},
			StartPoint:      kadmin.MostRecent,
			Filter: &kadmin.Filter{
				KeyFilter:       kadmin.StartsWithFilterType,
				KeySearchTerm:   "starts-with-key-term",
				ValueFilter:     kadmin.ContainsFilterType,
				ValueSearchTerm: "contains-value-term",
			},
			Limit: 500,
		}, &kadmin.ListedTopic{
			Name:           "topic1",
			PartitionCount: 10,
			Replicas:       3,
			RecordCount:    100,
		}, ui.NewTestKontext())

		// make sure form has been initialized
		m.View(ui.NewTestKontext(), ui.TestRenderer)

		render := m.View(ui.TestKontext, ui.TestRenderer)

		assert.Contains(t, render, "> Most Recent")
		assert.NotContains(t, render, "> Beginning")
		assert.Contains(t, render, "✓ 3")
		assert.Contains(t, render, "✓ 6")
		assert.Contains(t, render, "> ")
		assert.Contains(t, render, "starts-with-key-term")
		assert.Contains(t, render, "contains-value-term")
	})

	t.Run("submitting form loads consumption page with consumption information", func(t *testing.T) {
		m := New(&kadmin.ListedTopic{
			Name:           "topic1",
			PartitionCount: 10,
			Replicas:       1,
		}, ui.NewTestKontext())
		// make sure form has been initialized
		m.View(ui.NewTestKontext(), ui.TestRenderer)

		// select start from most recent
		cmd := m.Update(keys.Key(tea.KeyDown))
		cmd = m.Update(keys.Key(tea.KeyEnter))
		// next field
		m.Update(cmd())
		// select partition 3 and 5
		m.Update(keys.Key(tea.KeyDown))
		m.Update(keys.Key(tea.KeyDown))
		m.Update(keys.Key(tea.KeyDown))
		cmd = m.Update(keys.Key(tea.KeySpace))
		m.Update(keys.Key(tea.KeyDown))
		m.Update(keys.Key(tea.KeyDown))
		m.Update(keys.Key(tea.KeySpace))
		cmd = m.Update(keys.Key(tea.KeyEnter))
		// next field
		cmd = m.Update(cmd())
		// select limit 500
		m.Update(keys.Key(tea.KeyDown))
		cmd = m.Update(keys.Key(tea.KeyEnter))
		// next field
		cmd = m.Update(cmd())
		// next group
		m.Update(cmd())
		// no key filter
		cmd = m.Update(keys.Key(tea.KeyEnter))
		// next field
		cmd = m.Update(cmd())
		// no value filter
		msgs := keys.Submit(m)

		assert.Equal(t, nav.LoadConsumptionPageMsg{
			ReadDetails: kadmin.ReadDetails{
				TopicName: "topic1",
				Filter: &kadmin.Filter{
					KeySearchTerm:   "",
					ValueSearchTerm: "",
				},
				Limit:           500,
				PartitionToRead: []int{3, 5},
				StartPoint:      kadmin.MostRecent,
			},
			Topic: &kadmin.ListedTopic{
				Name:           "topic1",
				PartitionCount: 10,
				Replicas:       1,
				RecordCount:    0,
			},
		}, msgs[0])
	})

	t.Run("selecting partitions is optional", func(t *testing.T) {
		m := New(&kadmin.ListedTopic{
			Name:           "topic1",
			PartitionCount: 10,
			Replicas:       1,
		}, ui.NewTestKontext())
		// make sure form has been initialized
		m.View(ui.NewTestKontext(), ui.TestRenderer)

		// select start from most recent
		cmd := m.Update(keys.Key(tea.KeyDown))
		cmd = m.Update(keys.Key(tea.KeyEnter))
		// next field
		m.Update(cmd())
		// select no partitions
		cmd = m.Update(keys.Key(tea.KeyEnter))
		// next field
		cmd = m.Update(cmd())
		// select limit 500
		m.Update(keys.Key(tea.KeyDown))
		cmd = m.Update(keys.Key(tea.KeyEnter))
		// next field
		cmd = m.Update(cmd())
		// next group
		m.Update(cmd())
		// no key filter
		cmd = m.Update(keys.Key(tea.KeyEnter))
		// next field
		cmd = m.Update(cmd())
		// no value filter
		msgs := keys.Submit(m)

		assert.Equal(t, nav.LoadConsumptionPageMsg{
			ReadDetails: kadmin.ReadDetails{
				TopicName: "topic1",
				Filter: &kadmin.Filter{
					KeySearchTerm:   "",
					ValueSearchTerm: "",
				},
				Limit:           500,
				PartitionToRead: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
				StartPoint:      kadmin.MostRecent,
			},
			Topic: &kadmin.ListedTopic{
				Name:           "topic1",
				PartitionCount: 10,
				Replicas:       1,
				RecordCount:    0,
			},
		}, msgs[0])
	})

	t.Run("selecting key filter type starts-with displays key filter value field", func(t *testing.T) {
		m := New(&kadmin.ListedTopic{
			Name:           "topic1",
			PartitionCount: 10,
			Replicas:       1,
		}, ui.NewTestKontext())
		// make sure form has been initialized
		m.View(ui.NewTestKontext(), ui.TestRenderer)

		// select start from most recent
		cmd := m.Update(keys.Key(tea.KeyDown))
		cmd = m.Update(keys.Key(tea.KeyEnter))
		// next field
		m.Update(cmd())
		// select no partitions
		cmd = m.Update(keys.Key(tea.KeyEnter))
		// next field
		cmd = m.Update(cmd())
		// select limit 500
		m.Update(keys.Key(tea.KeyDown))
		cmd = m.Update(keys.Key(tea.KeyEnter))
		// next field
		cmd = m.Update(cmd())
		// next group
		m.Update(cmd())
		// starts-with key filter
		m.Update(keys.Key(tea.KeyDown))
		m.Update(keys.Key(tea.KeyDown))

		render := m.View(&kontext.ProgramKtx{
			Config:          nil,
			WindowWidth:     100,
			WindowHeight:    20,
			AvailableHeight: 20,
		}, ui.TestRenderer)

		assert.Contains(t, render, "Key Filter Term")

		t.Run("selecting no key filter type hides key value field again", func(t *testing.T) {
			m.Update(keys.Key(tea.KeyUp))
			m.Update(keys.Key(tea.KeyUp))

			render := m.View(&kontext.ProgramKtx{
				Config:          nil,
				WindowWidth:     100,
				WindowHeight:    20,
				AvailableHeight: 20,
			}, ui.TestRenderer)

			assert.NotContains(t, render, "Key Filter Value")
		})

		t.Run("selecting no key filter after filling in key filter term does not search for entered value", func(t *testing.T) {
			// select starts-with
			m.Update(keys.Key(tea.KeyDown))
			m.Update(keys.Key(tea.KeyDown))

			keys.UpdateKeys(m, "search-term")

			// selects none
			m.Update(keys.Key(tea.KeyUp))
			m.Update(keys.Key(tea.KeyUp))

			cmd = m.Update(keys.Key(tea.KeyEnter))
			// next field
			cmd = m.Update(cmd())
			// no value filter
			msgs := keys.Submit(m)

			assert.Equal(t, nav.LoadConsumptionPageMsg{
				ReadDetails: kadmin.ReadDetails{
					TopicName: "topic1",
					Filter: &kadmin.Filter{
						KeySearchTerm:   "",
						ValueSearchTerm: "",
					},
					Limit:           500,
					PartitionToRead: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
					StartPoint:      kadmin.MostRecent,
				},
				Topic: &kadmin.ListedTopic{
					Name:           "topic1",
					PartitionCount: 10,
					Replicas:       1,
					RecordCount:    0,
				},
			}, msgs[0])
		})
	})

	t.Run("filter on key value", func(t *testing.T) {
		m := New(&kadmin.ListedTopic{
			Name:           "topic1",
			PartitionCount: 10,
			Replicas:       1,
		}, ui.NewTestKontext())
		// make sure form has been initialized
		m.View(ui.NewTestKontext(), ui.TestRenderer)

		// select start from most recent
		cmd := m.Update(keys.Key(tea.KeyDown))
		cmd = m.Update(keys.Key(tea.KeyEnter))
		// next field
		m.Update(cmd())
		// select no partitions
		cmd = m.Update(keys.Key(tea.KeyEnter))
		// next field
		cmd = m.Update(cmd())
		// select limit 500
		m.Update(keys.Key(tea.KeyDown))
		cmd = m.Update(keys.Key(tea.KeyEnter))
		// next field
		cmd = m.Update(cmd())
		// next group
		m.Update(cmd())
		// starts-with key filter
		m.Update(keys.Key(tea.KeyDown))
		m.Update(keys.Key(tea.KeyDown))
		cmd = m.Update(keys.Key(tea.KeyEnter))
		// next field
		cmd = m.Update(cmd())
		// filter on key value search-term
		keys.UpdateKeys(m, "search-term")
		cmd = m.Update(keys.Key(tea.KeyEnter))
		// next field
		cmd = m.Update(cmd())
		// no value filter
		msgs := keys.Submit(m)

		assert.EqualValues(t, nav.LoadConsumptionPageMsg{
			ReadDetails: kadmin.ReadDetails{
				TopicName: "topic1",
				Filter: &kadmin.Filter{
					KeyFilter:       kadmin.StartsWithFilterType,
					KeySearchTerm:   "search-term",
					ValueSearchTerm: "",
				},
				Limit:           500,
				PartitionToRead: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
				StartPoint:      kadmin.MostRecent,
			},
			Topic: &kadmin.ListedTopic{
				Name:           "topic1",
				PartitionCount: 10,
				Replicas:       1,
				RecordCount:    0,
			},
		}, msgs[0])
	})

	t.Run("selecting value filter type starts-with displays filter value field", func(t *testing.T) {
		m := New(&kadmin.ListedTopic{
			Name:           "topic1",
			PartitionCount: 10,
			Replicas:       1,
		}, ui.NewTestKontext())
		// make sure form has been initialized
		m.View(ui.NewTestKontext(), ui.TestRenderer)

		// select start from most recent
		cmd := m.Update(keys.Key(tea.KeyDown))
		cmd = m.Update(keys.Key(tea.KeyEnter))
		// next field
		m.Update(cmd())
		// select no partitions
		cmd = m.Update(keys.Key(tea.KeyEnter))
		// next field
		cmd = m.Update(cmd())
		// select limit 500
		m.Update(keys.Key(tea.KeyDown))
		cmd = m.Update(keys.Key(tea.KeyEnter))
		// next field
		cmd = m.Update(cmd())
		// next group
		m.Update(cmd())
		// no key filter
		cmd = m.Update(keys.Key(tea.KeyEnter))
		// next field
		cmd = m.Update(cmd())
		// starts-with value filter
		m.Update(keys.Key(tea.KeyDown))
		m.Update(keys.Key(tea.KeyDown))

		render := m.View(&kontext.ProgramKtx{
			Config:          nil,
			WindowWidth:     100,
			WindowHeight:    20,
			AvailableHeight: 20,
		}, ui.TestRenderer)

		assert.Contains(t, render, "Value Filter Term")

		// make sure the value filter term field is focussed
		cmd = m.Update(keys.Key(tea.KeyEnter))
		// next field
		cmd = m.Update(cmd())
		// next group
		m.Update(cmd())

		field := m.form.GetFocusedField()
		assert.IsType(t, &huh.Input{}, field)
		assert.Contains(t, field.View(), "Value Filter Term")

		t.Run("selecting no value filter type hides key filter value field again", func(t *testing.T) {
			cmd = m.Update(keys.Key(tea.KeyShiftTab))
			// prev field
			cmd = m.Update(cmd())

			m.Update(keys.Key(tea.KeyUp))
			m.Update(keys.Key(tea.KeyUp))

			render := m.View(&kontext.ProgramKtx{
				Config:          nil,
				WindowWidth:     100,
				WindowHeight:    20,
				AvailableHeight: 20,
			}, ui.TestRenderer)

			assert.NotContains(t, render, "Value Filter Term")
		})

		t.Run("selecting no value filter after filling in a value filter term does not search for entered value", func(t *testing.T) {
			// select starts-with
			m.Update(keys.Key(tea.KeyDown))
			m.Update(keys.Key(tea.KeyDown))

			keys.UpdateKeys(m, "search-term")

			// selects none
			m.Update(keys.Key(tea.KeyUp))
			m.Update(keys.Key(tea.KeyUp))

			cmd = m.Update(keys.Key(tea.KeyEnter))
			// next field
			cmd = m.Update(cmd())
			// no value filter
			msgs := keys.Submit(m)

			assert.Equal(t, nav.LoadConsumptionPageMsg{
				ReadDetails: kadmin.ReadDetails{
					TopicName: "topic1",
					Filter: &kadmin.Filter{
						KeySearchTerm:   "",
						ValueSearchTerm: "",
					},
					Limit:           500,
					PartitionToRead: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
					StartPoint:      kadmin.MostRecent,
				},
				Topic: &kadmin.ListedTopic{
					Name:           "topic1",
					PartitionCount: 10,
					Replicas:       1,
					RecordCount:    0,
				},
			},
				msgs[0])
		})
	})
}
