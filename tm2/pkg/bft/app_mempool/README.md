# Application-Level Mempool (app_mempool)

This folder contains an experimental implementation of an **application-aware mempool**, designed as a future alternative to the current FIFO-based `CListMempool` used in Gnoland.

---

## ✍️ Background

This implementation is based on the RFC:

**📄 [Simplified Mempool Design](https://github.com/gnolang/gno/issues/4257)**  
Authors: Miloš Živković (@zivkovicmilos), Dragan Milošević (@Milosevic02)  
Date: May 7, 2025  
Status: Draft

The RFC addresses several UX and performance issues in the current Tendermint2 mempool design, including:

- Out-of-order transaction inclusion  
- Lack of nonce-awareness  
- Lack of transaction prioritization based on gas fees  
- Inefficient locking and throughput bottlenecks  

---

## ⚙️ What this implements

This `app_mempool` version includes:

- Per-account mempool state with separation of `pending` (executable) and `queued` (future) transactions  
- Strict transaction ordering by nonce  
- Optional prioritization by `Price()`  
- Transaction promotion logic as gaps are filled  
- Lightweight concurrency via per-account locks  
- Full support for `Flush`, `Update`, `Content`, and targeted access methods (`Tx`, `Size`, etc.)

All logic is implemented using a custom `Tx` interface and backed by a high-performance insertion queue:

```go
type Tx interface {
	Hash() []byte
	Sender() crypto.Address
	Sequence() uint64
	Price() uint64
	Size() uint64
}
````

---

## ⚠️ Why it's not active (yet)

This implementation assumes that each transaction exposes both the **sender** and **sequence number (nonce)**.
However, in the current Gno setup, the ABCI layer **does not expose** this metadata at the time the transaction is received.

This means:

* We cannot extract `Sender()` or `Sequence()` directly from the raw transaction.
* Therefore, the proposed logic cannot reliably function in the real node (yet).

---

## 📁 Related work

If the required metadata becomes available in the future, this design becomes a highly attractive replacement.
Until then, a minimal, cleaner implementation is available for use:

👉 [Simplified Mempool](https://github.com/Milosevic02/gno/blob/feat/mempool/tm2/pkg/bft/my_mempool/mempool.go) – a lightweight, FIFO-based alternative to CListMempool, featuring a simplified interface, lower complexity, and significantly improved performance characteristics.

---

## ✅ Summary

Although currently inactive, this `app_mempool` stands as a complete and working implementation of the RFC logic, awaiting support from the underlying infrastructure.

It demonstrates:

* A cleaner, more structured approach to transaction handling
* Readiness to support fee markets, DoS protection, and higher throughput
* Clear separation between application logic and transaction flow

The code remains available as a reference, benchmark baseline, and potential future integration candidate.

