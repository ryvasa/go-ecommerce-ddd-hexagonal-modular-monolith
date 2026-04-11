---
trigger: always_on
---

# AI Agent Rules — Backend Golang

## Tujuan

Dokumen ini berisi **aturan keras (non-negotiable rules)** untuk AI agent saat membantu pengembangan backend **Golang** menggunakan **DDD (Domain-Driven Design)** dan **Hexagonal Architecture (Ports & Adapters)**. Fokus utama: _maintainability_, _testability_, _explicit boundaries_, dan _long-term scalability_.

---

## 1. Prinsip Umum (Hard Principles)

1. **Domain adalah pusat segalanya**
   - Domain **tidak boleh** bergantung pada framework, database, HTTP, gRPC, ORM, atau library eksternal.
   - Domain hanya berisi _business rules_.

2. **Dependency Rule (Clean Architecture)**
   - Arah dependensi selalu **ke dalam**:
     `Delivery → Application → Domain`
   - Domain **tidak tahu** siapa yang memanggilnya.

3. **Explicit over Implicit**
   - Tidak ada magic, side-effect tersembunyi, atau behavior implisit.
   - Semua dependency harus jelas (constructor, interface, atau DI).

4. **Fail Fast & Deterministic**
   - Error harus muncul sedini mungkin.
   - Jangan menyembunyikan error atau mengandalkan default behavior.

---

## 2. Struktur Folder (Wajib)

```
cmd/
  api/
    main.go
internal/
  <module>/
    domain/
      entity/
      valueobject/
      repository.go
      service.go (jika domain service diperlukan)
      error.go
    application/
      usecase/
      dto/
      port/
        in/
        out/
    infrastructure/
      persistence/
      http/
      grpc/
      messaging/
    delivery/
      http/
      grpc/
pkg/
  config/
  logger/
  db/
```

> AI **dilarang** mencampur layer (misalnya HTTP handler langsung mengakses repository).

---

## 3. Domain Layer Rules

### 3.1 Entity

- Entity **selalu punya identity**.
- Tidak boleh ada setter publik sembarangan.
- State hanya berubah lewat **method domain**.

```go
func (u *User) Activate() error {
    if u.status == StatusActive {
        return ErrAlreadyActive
    }
    u.status = StatusActive
    return nil
}
```

### 3.2 Value Object

- Immutable.
- Validasi di constructor.
- Tidak expose primitive mentah.

```go
func NewEmail(v string) (Email, error)
```

### 3.3 Domain Error

- Error domain **bermakna bisnis**, bukan teknis.
- Tidak menggunakan error string mentah.

---

## 4. Application Layer Rules

### 4.1 Usecase

- Usecase **1 tanggung jawab**.
- Tidak ada logika HTTP / DB / gRPC.
- Semua dependency lewat **port/out**.

```go
type RegisterUserUsecase struct {
    userRepo port.UserRepository
}
```

### 4.2 DTO

- DTO hanya untuk **boundary crossing**.
- DTO **tidak** dipakai di domain.

### 4.3 Transaction Handling

- Transaction dikontrol di **application layer**.
- Domain tidak tahu transaction.

---

## 5. Port Rules (Hexagonal)

### 5.1 Port In

- Mewakili _intent bisnis_.
- Biasanya interface usecase.

### 5.2 Port Out

- Didefinisikan oleh **caller module**, bukan implementor.
- Domain/application hanya bergantung pada interface.

```go
type UserRepository interface {
    Save(ctx context.Context, user *User) error
}
```

---

## 6. Infrastructure Layer Rules

1. Implementasi detail teknis saja.
2. ORM, HTTP client, cache, queue **hanya di sini**.
3. Mapping entity ↔ persistence model dilakukan di layer ini.

> AI **dilarang** memunculkan GORM / SQL di domain atau usecase.

---

## 7. Delivery Layer Rules

1. Delivery hanya:
   - parsing request
   - validasi input ringan
   - mapping ke DTO
   - memanggil usecase
2. Tidak ada business logic.

---

## 8. Dependency Injection (DI)

1. Semua wiring dilakukan di **composition root** (cmd/main).
2. Domain dan application **tidak pernah** tahu DI tool.
3. Gunakan constructor injection.

---

## 9. Testing Rules

1. Domain: pure unit test, tanpa mock infra.
2. Usecase: mock port/out.
3. Infrastructure: integration test.
4. Jangan mock domain entity.

---

## 10. Anti-Patterns (DILARANG)

❌ Fat service / God object
❌ Entity + ORM annotation campur
❌ Repository return DTO
❌ Usecase return HTTP response
❌ Domain memanggil repository langsung
❌ Cross-module import tanpa port

---

## 11. AI Agent Decision Checklist

Sebelum menulis atau mengubah kode, AI **WAJIB** memastikan:

### 11.1 Boundary & Layer

- [ ] Kode ini berada di **module apa** (bounded context mana)?
- [ ] Apakah module ini pemilik domain tersebut?
- [ ] Kode ini berada di **layer yang benar** (domain / application / adapter / shared)?
- [ ] Apakah ada import yang melanggar batas module?
- [ ] Apakah arah dependensi selalu menuju ke dalam?

### 11.2 Domain Purity

- [ ] Apakah ini benar-benar _business rule_?
- [ ] Apakah domain bebas dari ORM (GORM), SQL, HTTP, JSON tag, config?
- [ ] Apakah entity memiliki behavior (tidak anemic)?
- [ ] Apakah invariant dijaga di constructor atau method domain?
- [ ] Apakah value object immutable dan tervalidasi?

### 11.3 Repository & Port

- [ ] Apakah repository interface berada di **domain**?
- [ ] Apakah repository menggunakan entity & value object (bukan DTO)?
- [ ] Apakah cross-module call dilakukan lewat **port**, bukan import langsung?
- [ ] Apakah module pemilik domain yang mengimplementasikan port tersebut?

### 11.4 Application Layer

- [ ] Apakah usecase hanya mengorkestrasi flow?
- [ ] Apakah usecase bebas dari SQL, ORM, HTTP, dan framework?
- [ ] Apakah transaction dikelola di application layer?
- [ ] Apakah semua dependency di-inject via constructor?

### 11.5 Adapter / Infrastructure

- [ ] Apakah GORM hanya digunakan di adapter/out/persistence?
- [ ] Apakah entity domain bebas dari tag GORM?
- [ ] Apakah mapping entity ↔ persistence model eksplisit?
- [ ] Apakah Casbin hanya digunakan di adapter / middleware?

### 11.6 Tooling & Ops

- [ ] Apakah Wire hanya digunakan di composition root (cmd)?
- [ ] Apakah tidak ada auto-migrate di production?
- [ ] Apakah migration versioned dan reproducible?
- [ ] Apakah config tidak bocor ke domain?

### 11.7 Testability

- [ ] Apakah domain bisa diuji tanpa database atau infra?
- [ ] Apakah usecase bisa diuji dengan mock port?
- [ ] Apakah adapter diuji via integration test?
- [ ] Apakah perubahan ini tidak memaksa integration test untuk domain?

Jika **satu saja** checklist di atas gagal → AI **HARUS berhenti**, menjelaskan pelanggaran, dan memperbaiki desain sebelum menulis kode.

---

## Penutup

Jika ada konflik antara _kecepatan_ dan _arsitektur_, **arsitektur menang**.
AI agent harus selalu memilih solusi yang paling eksplisit, terpisah, dan mudah diuji.

---

# Tambahan Aturan — Modular Monolith (DDD-heavy + Hexagonal)

Dokumen ini memperketat rule sebelumnya dengan fokus **Modular Monolith**, serta standar tooling: **Wire**, **GORM**, **Database Migration**, dan **Casbin (industry standard)**.

---

## 12. Prinsip Modular Monolith (WAJIB)

1. **Module = Bounded Context**
   - Satu module merepresentasikan satu domain bisnis.
   - Module bukan sekadar folder.

2. **No Cross-Domain Import**
   - `domain` suatu module **TIDAK BOLEH** di-import module lain.
   - Komunikasi antar module hanya lewat **application port/out**.

3. **Explicit Module Contract**
   - Jika module A butuh data module B:
     - A mendefinisikan port/out
     - B mengimplementasikan adapter/out
     - Wiring hanya di composition root

4. **Shared Tidak Mengandung Business Rule**
   - `shared` hanya berisi:
     - logging
     - transaction manager
     - event bus
     - middleware
   - Tidak boleh ada entity, usecase, atau repository di shared.

---

## 13. Domain Rules (DDD-Heavy)

1. **Domain adalah pemilik kebenaran bisnis**
   - Semua invariant, rule, dan state transition ada di domain.

2. **Repository Interface WAJIB di Domain**
   - Repository adalah konsep bisnis (collection of aggregates).
   - Repository **bukan** detail teknis.

3. **Domain Tidak Mengenal Usecase**
   - Domain tidak tahu siapa yang memanggil.
   - Domain tidak mengorkestrasi flow.

4. **Entity Tidak Anemic**
   - Entity wajib punya behavior.
   - Tidak boleh hanya struct + getter/setter.

---

## 14. Application Layer Rules

1. **Application = Orchestrator**
   - Mengatur flow antar entity dan module.
   - Tidak menyimpan state bisnis.

2. **Usecase Tidak Mengandung SQL / ORM**
   - Semua persistence lewat repository interface.

3. **Cross-Module Call = Port Only**
   - Tidak ada direct import domain module lain.

4. **Transaction Dikelola di Application**
   - Menggunakan shared transaction manager.

---

## 15. Adapter / Infrastructure Rules

1. **Adapter = Detail Teknis**
   - HTTP, DB, Cache, Message Broker, Security.

2. **GORM HANYA di Adapter Out**
   - Tidak boleh muncul di domain / application.
   - Mapping entity ↔ model wajib eksplisit.

3. **Model Persistence Terpisah dari Entity**
   - Entity tidak punya tag GORM.

---

## 16. Dependency Injection (Wire)

1. **Wire HANYA di Composition Root**
   - `cmd/<app>/wire.go`

2. **Domain Tidak Tahu Wire**
   - Application juga tidak tahu Wire.

3. **Semua Dependency via Constructor**

---

## 17. Authorization (Casbin – Industry Standard)

1. **Casbin BUKAN Domain Concern**
   - Authorization adalah infrastructure concern.

2. **Policy Tidak Berisi Business Logic**
   - Policy hanya rule akses (who can do what).

3. **Domain Tidak Mengecek Role / Permission**
   - Domain mengecek state bisnis, bukan otorisasi.

---

## 18. Database & Migration

1. **Migration Bersifat Global**
   - Disimpan di root `migrations/`.

2. **Schema Mewakili Domain, Bukan ORM**
   - Nama tabel & kolom mencerminkan bahasa bisnis.

3. **No Auto-Migrate di Production**
   - Migration harus eksplisit dan versioned.

---

## 19. Testing Strategy (Modular Monolith)

1. **Domain Test = Pure Unit Test**
   - Tanpa mock infra.

2. **Application Test = Mock Port**

3. **Adapter Test = Integration Test**
   - DB sungguhan / test container.

4. **Cross-Module Contract Test**
   - Pastikan port/out tidak berubah sembarangan.

---

## 20. Absolute Anti-Patterns (NON-NEGOTIABLE)

❌ Cross import domain antar module
❌ Entity dengan tag GORM
❌ Usecase return HTTP response
❌ Shared domain model
❌ Domain melakukan authorization
❌ Auto-migrate tanpa kontrol

---

## Prinsip Penutup (WAJIB DIPEGANG AI)

> **Modular Monolith bukan Microservice versi hemat.**
>
> Ia adalah monolith dengan _domain discipline_.
>
> Jika ragu antara cepat vs benar → **pilih benar**.
