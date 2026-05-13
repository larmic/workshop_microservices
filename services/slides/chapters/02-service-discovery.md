## Beispiel: Service Discovery

```go [1-2|3-6]
resolver := consul.NewResolver("flight-service")
addr, err := resolver.Resolve(ctx)
if err != nil {
    return fmt.Errorf("resolve: %w", err)
}
resp, err := http.Get(addr + "/flights")
```

---

## Zur Hands-on-Story

<div class="box">

**Story 2 — Service Discovery**

Im Dashboard öffnen und der Anleitung folgen.

</div>
