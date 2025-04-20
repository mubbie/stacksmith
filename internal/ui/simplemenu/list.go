// ui/simplemenu/list.go
package simplemenu

import (
    "fmt"
    "github.com/mubbie/stacksmith/internal/ui/styles"
)

// SelectableList provides a reusable list component
type SelectableList struct {
    Items     []string
    Cursor    int
    Selected  map[int]bool
    OrderMap  map[string]int // Maps item value to selection order
    NextOrder int            // Next order index to assign 
}

// NewSelectableList creates a new selectable list
func NewSelectableList(items []string) *SelectableList {
    return &SelectableList{
        Items:     items,
        Selected:  make(map[int]bool),
        OrderMap:  make(map[string]int),
        NextOrder: 1,
    }
}

// MoveUp moves the cursor up
func (l *SelectableList) MoveUp() {
    if l.Cursor > 0 {
        l.Cursor--
    }
}

// MoveDown moves the cursor down
func (l *SelectableList) MoveDown() {
    if l.Cursor < len(l.Items)-1 {
        l.Cursor++
    }
}

// ToggleSelected toggles selection for current item
func (l *SelectableList) ToggleSelected() {
    current := l.Cursor
    currentItem := l.Items[current]
    
    if l.Selected[current] {
        // Deselect the item
        l.Selected[current] = false
        delete(l.OrderMap, currentItem)
        
        // Reorder the remaining items
        for item, order := range l.OrderMap {
            if order > l.OrderMap[currentItem] {
                l.OrderMap[item]--
            }
        }
        l.NextOrder--
    } else {
        // Select the item
        l.Selected[current] = true
        l.OrderMap[currentItem] = l.NextOrder
        l.NextOrder++
    }
}

// GetSelectedItems returns selected items in order of selection
func (l *SelectableList) GetSelectedItems() []string {
    if len(l.Selected) == 0 {
        return nil
    }
    
    result := make([]string, len(l.OrderMap))
    
    // Map order -> item
    orderToItem := make(map[int]string)
    for item, order := range l.OrderMap {
        orderToItem[order] = item
    }
    
    // Add items in order
    for i := 1; i <= len(l.OrderMap); i++ {
        if item, ok := orderToItem[i]; ok {
            result[i-1] = item
        }
    }
    
    return result
}

// GetSelectedCount returns the number of selected items
func (l *SelectableList) GetSelectedCount() int {
    return len(l.OrderMap)
}

// Render renders the list with selection indicators
func (l *SelectableList) Render(showCheckboxes bool, showOrder bool) string {
    result := ""
    
    for i, item := range l.Items {
        cursor := styles.CursorStyle(i == l.Cursor)
        
        itemStyle := styles.Normal
        if i == l.Cursor {
            itemStyle = styles.Selected
        }
        
        line := cursor + " "
        
        if showCheckboxes {
            checkbox := "[ ]"
            if l.Selected[i] {
                checkbox = "[x]"
            }
            line += checkbox + " "
        }
        
        // Add order number if requested and item is selected
        if showOrder && l.Selected[i] {
            if order, ok := l.OrderMap[item]; ok {
                line += fmt.Sprintf(" %d ", order)
            } else {
                line += "   "
            }
        } else if showCheckboxes {
            line += "   " // Maintain spacing
        }
        
        line += itemStyle.Render(item)
        result += line + "\n"
    }
    
    return result
}