# tmhash – Hash Pool Optimization Attempt

This module explores an optimization for reducing memory allocations during SHA256 hashing using `sync.Pool`.

## Idea

In performance-critical paths like the mempool, frequent hash operations can lead to many short-lived allocations.  
The idea was to reuse `sha256.Hash` instances via a pool to reduce GC pressure and improve performance.

## Implementation

- Introduced `hashPool` and `truncatedHashPool` for full and truncated SHA256 hashers.
- All functions (`New`, `Sum`, `NewTruncated`, `SumTruncated`) preserve the existing API.
- Instances are reset and returned to the pool after use.

## Result

Despite correct functionality, benchmarks showed **slightly worse performance** due to heap usage and pool overhead.  
In this case, Go’s native stack allocation for short-lived hashers is more efficient.

This experiment is not part of the final implementation but remains here as a reference.

