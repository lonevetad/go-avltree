# In-order link repair plan for the AVL tree

## 1. What the current rebuild helpers do

The tree currently has two different kinds of auxiliary links:

- `prevInOrder` / `nextInOrder`: represent the node order as a circular doubly linked list in in-order traversal order.
- `prevInserted` / `nextInserted`: represent chronological insertion order.

The functions below are the current safety net for the in-order links.

### `rebuildMetadataAndInOrderLinks`

This helper recursively walks the subtree rooted at `node` and:

1. resets the parent pointer for the current node,
2. recursively rebuilds the metadata for the left and right subtrees,
3. recomputes:
   - `sizeLeft`
   - `sizeRight`
   - `height`
4. re-establishes the parent links for each child.

Its purpose is mainly structural correctness of the tree itself: sizes and heights must match the actual subtree shape after rotations or removals.

### `rebuildInOrderLinks`

This helper rebuilds the in-order traversal links from scratch:

1. it calls `rebuildMetadataAndInOrderLinks` first,
2. it performs an in-order DFS over the tree,
3. it reconnects each node to its predecessor and successor in the traversal order,
4. it closes the circular list by making the last node point back to the first one,
5. it updates `minValue` to the first node in the in-order list.

This is the function that restores the correctness of the traversal-linked-list view of the tree.

---

## 2. Why the links can break during removal

The AVL tree is not just a binary tree; it is also expected to expose an ordered iteration view through the in-order links.

When a node is removed, the following operations can change the structure:

- the node itself is removed or replaced,
- a child is promoted,
- a successor is swapped into place,
- rotations are performed during rebalancing,
- parent/child relationships are changed.

If the in-order links are not updated at the same time, the tree can still be structurally balanced, but traversal becomes invalid:

- a node may point to the wrong predecessor/successor,
- the list may be broken into two fragments,
- `minValue` may no longer point to the true first node,
- iteration can skip nodes or loop incorrectly.

That is why the current implementation calls `rebuildInOrderLinks()` at the end of `remove()`.

---

## 3. Why the current approach works but is expensive

The current approach is simple and reliable:

- it guarantees correctness,
- it avoids subtle link bugs,
- it works even when the tree shape is complex after rotations.

However, it scans the whole tree, so the removal operation becomes $O(n)$ in the worst case.

The goal of the new implementation should be:

- maintain correctness,
- repair the links incrementally during the delete operation,
- keep the overall complexity close to $O(\log n)$ for a balanced AVL tree.

---

## 4. Desired invariants for the in-order links

After every successful deletion, the following invariants should hold:

1. The tree is still a valid AVL tree.
2. The in-order traversal order is represented by the `prevInOrder` / `nextInOrder` links.
3. The list is circular and contains exactly the live nodes.
4. The node pointed to by `minValue` is the first node in the in-order order.
5. The sentinel node (`_NIL`) has correct list links that do not accidentally point into the live tree.

These are the properties that must be preserved while editing links.

---

## 5. Plan for fixing the links during removal

The fix should be done locally, not by rebuilding the entire tree at the end.

### Phase A: identify the affected neighborhood

For each removal case, determine:

- the node being removed,
- its predecessor in in-order order,
- its successor in in-order order,
- the parent of the node or the node that will replace it,
- whether the removed node was the current minimum.

This gives enough context to update the local list neighborhood.

### Phase B: repair the in-order links for each deletion case

#### Case 1: removing a leaf

If the node has no children:

- splice it out of the circular list by connecting its predecessor and successor directly,
- if the removed node was the minimum, update `minValue` to its successor,
- if the removed node was the root and the tree becomes empty, reset the list to the sentinel state.

Pseudo-logic:

```go
pred := n.prevInOrder
succ := n.nextInOrder

if pred != t._NIL {
    pred.nextInOrder = succ
} else {
    // n was the first node in the list
}

if succ != t._NIL {
    succ.prevInOrder = pred
} else {
    // n was the last node in the list
}
```

#### Case 2: removing a node with one child

If a child replaces the removed node:

- the child must inherit the predecessor and successor of the removed node,
- the predecessor must point to the child,
- the successor must point back to the child,
- if the removed node was the minimum, update `minValue` to the child.

This is the most common local repair case.

#### Case 3: removing a node with two children

This is the most subtle case.

The usual AVL delete pattern is:

- find the in-order successor,
- copy its key/value into the node being removed,
- delete the successor instead.

That means the real structural removal happens on the successor, but the visible logical node remains in place.

In this case:

- the successor is the node that will actually leave the tree,
- the predecessor/successor of that successor are the ones that must be spliced together,
- the node that remains in place must keep the correct in-order neighbors.

The repair logic should therefore be applied to the actual node being unlinked from the tree, while preserving the logical position of the replacement node.

### Phase C: repair the metadata after the local link change

Once the node has been unlinked and the list neighborhood is correct, the tree metadata should be updated only on the affected ancestors.

That means reusing the existing helpers such as:

- `recalculateHeight(...)`
- `recalculateSizes(...)`
- `insertFixup(...)`

The important point is that these updates should happen after the local link repair, not as a full-tree rebuild.

### Phase D: restore sentinel links

At the end of the operation, ensure:

- `_NIL.prevInOrder` and `_NIL.nextInOrder` are correct,
- `_NIL` does not accidentally appear inside the live list,
- the circular list remains well formed.

---

## 6. Suggested helper functions

To keep the implementation readable, the fix can be split into small helpers:

- `spliceOutFromInOrderList(node)`
- `replaceInOrderNode(oldNode, newNode)`
- `linkInOrderNeighbors(pred, succ)`
- `updateMinValueAfterRemoval(removedNode, replacementNode)`

These helpers should be responsible only for the in-order links and the minimum pointer.

The metadata and AVL balancing logic can stay in the existing routines.

---

## 7. Alternative fixes

### Alternative A: keep the full rebuild

This is the current approach.

Pros:

- simplest,
- less bug-prone,
- easy to reason about.

Cons:

- $O(n)$ delete time,
- rebuilds a lot of data that may not be necessary.

### Alternative B: local repair with explicit pointer updates

This is the recommended approach.

Pros:

- preserves $O(\log n)$ behavior for a balanced AVL tree,
- keeps the logic localized,
- avoids unnecessary full-tree traversal.

Cons:

- more error-prone,
- requires careful handling of the three deletion cases,
- must preserve all invariants manually.

### Alternative C: maintain a separate ordered container

One could maintain a secondary ordered structure such as a linked list or an order-statistics structure that is updated together with the tree.

Pros:

- simpler to reason about traversal order,
- potentially easier to keep in sync.

Cons:

- duplicates state,
- increases maintenance cost,
- can become inconsistent if one structure is updated without the other.

### Alternative D: rebuild only when the tree becomes inconsistent

This is a hybrid strategy.

Pros:

- fast in the common case,
- keeps complexity low on average.

Cons:

- harder to prove correctness,
- may silently break if the inconsistency is not detected.

---

## 8. Recommended implementation strategy

The best path is:

1. keep the current AVL balancing logic,
2. add local helper functions for in-order link repair,
3. call those helpers in each removal branch,
4. update the minimum pointer and sentinel links explicitly,
5. keep the full rebuild only as a fallback or debug assertion if needed.

This gives a practical balance between correctness, clarity, and performance.

---

## 9. Practical note

The important thing is not to treat the in-order links as an afterthought. They are part of the tree’s public behavior, not just an internal implementation detail. If iteration depends on them, then every deletion must update them in a way that preserves the exact traversal order.

In other words: the tree must remain correct both as an AVL tree and as an ordered linked structure.
