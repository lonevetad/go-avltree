# AVL tree test and implementation plan

## Goals

- Strengthen the test suite around the public API instead of only low-level rotations.
- Reuse shared helpers for creating trees and data so new scenarios are easy to add.
- Make remove-by-key behavior easier to reason about and cover with table-driven tests.

## Short-term improvements

1. Keep the current helper-based construction pattern for test data.
   - Reuse newTreeWithKeys for inserting a known sequence of keys.
   - Add small helpers for common scenarios such as "single node", "two-node tree", and "balanced tree".
2. Expand the remove-by-key suite with table-driven cases.
   - Cover removing a leaf.
   - Cover removing a missing key.
   - Add cases for removing the root and removing internal nodes once the implementation is stable enough.
3. Add public API assertions around traversal and indexing.
   - Verify GetAt returns keys in sorted order.
   - Verify ForEach for InOrder, ReverseInOrder, Queue, and Stack.
   - Verify tree size and root/min-value state after each mutation.
4. Add low-level structural assertions for remove-by-key cases.
   - Verify minValue and firstInserted still point at the correct nodes after a deletion.
   - Verify in-order and chronological traversals remain consistent with the expected key sequence.
   - Reuse the metadata checker to ensure cached heights and subtree sizes remain correct.
5. Separate regression tests from long-running structural tests.
   - Keep large integration tests like TestPrint_NTT but mark them as slower or optional if needed.
   - Add focused tests that fail fast and describe a single behavior.

## Medium-term improvements

1. Investigate remove-by-key correctness for more complex trees.
   - Reproduce the cases that delete a node with one child or two children.
   - Compare the resulting shape, ordering, and cached metadata to an expected reference tree.
2. Add test helpers for expected tree shape.
   - A helper could build a reference tree from explicit keys and compare it to the implementation.
3. Improve failure messages.
   - Make assertions print the tree shape and the key order after a failed remove, so debugging is faster.
4. Consider splitting the giant test file.
   - Keep rotation tests in one file, public API tests in another, and remove tests in a third.

## Suggested next cases to add

- Remove a node with one child.
- Remove the root from a two-node tree.
- Remove an internal node from a larger tree and ensure the in-order sequence remains correct.
- Remove repeatedly from the same tree and verify both size and ordering after each step.

## Long-running structural test tackling

The existing structural test, TestPrint_NTT, is meant to stress the tree over many insertions and check that its internal ordering, link structure, and cached metadata remain consistent across a broader range of shapes. It is intentionally heavier than the focused regression tests, so it is useful for catching subtle balance or pointer issues, but it is also slower and noisier when it fails. The plan is to keep it as a broader verification test while adding smaller, faster regression cases that pinpoint the exact behavior being exercised.
