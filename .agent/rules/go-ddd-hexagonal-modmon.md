---
trigger: always_on
---

AI Agent Rules — Backend Golang (DDD & Hexagonal)

1. Prinsip Utama (Hard Principles)
   Domain Centric: Domain tidak boleh bergantung pada library, framework, atau external tools. Hanya berisi business logic.

Dependency Rule: Arah dependensi selalu ke dalam: Delivery -> Application -> Domain.

Explicit over Implicit: Gunakan constructor, interface, dan DI. Hindari magic atau side-effects.

Fail Fast: Error ditangani sedini mungkin secara deterministik.

Modular Monolith: Satu modul = satu Bounded Context. Dilarang keras cross-import antar domain modul.

2. Struktur Folder Wajib
   Plaintext
   internal/<module>/
   domain/ # Entity, ValueObject, Repository Interface, Domain Service
   application/ # Usecase (Orchestrator), DTO, Port (In/Out Interface)
   infrastructure/ # Persistence (GORM), External Clients, Messaging
   delivery/ # HTTP/gRPC Handlers
   pkg/ # Shared libraries (non-business logic)
   cmd/ # Main entry & Composition Root (Wire)
3. Layer Rules
   3.1 Domain Layer
   Entity: Punya identitas, state hanya berubah via domain method (bukan public setter).

Value Object: Immutable, divalidasi di constructor.

Repository Interface: Wajib di layer ini (abstraksi bisnis).

No Anemic Model: Entity harus memiliki behavior, bukan hanya sekadar struct data.

3.2 Application Layer
Usecase: Satu tanggung jawab. Tidak boleh ada logika HTTP/DB.

DTO: Hanya untuk boundary crossing, tidak masuk ke layer domain.

Transaction: Dikelola di sini menggunakan Shared Transaction Manager.

Ports: Port In (interface usecase), Port Out (interface yang dibutuhkan usecase).

3.3 Infrastructure & Delivery
Persistence: GORM/SQL dilarang keras muncul di domain/usecase. Mapping Entity ↔ Persistence Model harus eksplisit.

Delivery: Hanya fokus parsing, mapping DTO, dan memanggil usecase.

Casbin: Authorization adalah infrastruktur/middleware, bukan domain concern.

4. Testing Strategy (Pyramid)
   4.1 Unit Test (Fast & Cheap)
   Domain: Wajib test business rules & invariants. Dilarang Mock.

Usecase: Test orkestrasi dengan Mock Port Out.

Pattern: Gunakan Table-driven tests.

4.2 Integration & E2E
Integration: Wajib menggunakan Real DB (Testcontainers/Docker) untuk Repository & Adapter.

E2E: Test full user flow (HTTP Request -> DB). Simpan di /test/e2e/.

Build Tags: Gunakan //go:build integration atau //go:build e2e.

5. Dependency Injection (Wire)
   Wiring hanya dilakukan di Composition Root (cmd/).

Domain & Application dilarang tahu keberadaan DI Tool.

Gunakan Constructor Injection.

6. Anti-Patterns (DILARANG)
   ❌ Entity memiliki tag GORM/JSON (kecuali di layer adapter).

❌ Repository mengembalikan DTO (harus Entity).

❌ Usecase mengembalikan HTTP response status code.

❌ Cross-domain import tanpa port/out.

❌ Auto-migrate di production (wajib versioned migration).

7. AI Agent Decision Checklist
   AI WAJIB berhenti dan memperbaiki desain jika salah satu poin di bawah dilanggar:

Boundary: Apakah kode berada di module dan layer yang benar?

Purity: Apakah domain bebas dari GORM, SQL, atau library eksternal?

Port: Apakah komunikasi antar-modul melalui interface port/out?

Mocking: Apakah unit test domain menggunakan mock? (Harusnya tidak).

Logic Leak: Apakah business logic bocor ke delivery atau infrastructure?

Penutup: Jika ragu antara kecepatan dan arsitektur, Arsitektur Menang. AI harus memilih solusi yang paling eksplisit dan testable.
