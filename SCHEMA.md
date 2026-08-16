# Agriculture Shop — Data Schema (v2)

Finalized schema for the POS system. camelCase field names throughout.
No Suppliers/Purchase Orders, no Stock Movement Log, no Employees (merged into Users).

## Collections

### users
| field | type | note |
|---|---|---|
| id | ObjectId | |
| username | string | unique |
| passwordHash | string | |
| role | string | ADMIN / CASHIER / STOCK |
| firstName | string | |
| lastName | string | |
| phone | string | |
| status | string | ACTIVE / INACTIVE |

### units
| field | type | note |
|---|---|---|
| id | ObjectId | |
| name | string | ອັນ, ແກ້ວ, ຖົງ, ລັງ, ກິໂລ |

### categories
| field | type | note |
|---|---|---|
| id | ObjectId | |
| name | string | |

### products
| field | type | note |
|---|---|---|
| id | ObjectId | |
| barcode | string | unique, optional |
| name | string | |
| categoryId | ObjectId | ref categories |
| unitId | ObjectId | ref units |
| costPrice | float | |
| retailPrice | float | |
| wholesalePrice | float | |
| wholesaleMinQty | int | qty >= this uses wholesalePrice, per-product |
| stockQty | float | |
| minStockAlert | float | |
| status | string | ACTIVE / DISCONTINUED |

### customers
| field | type | note |
|---|---|---|
| id | ObjectId | |
| name | string | |
| phone | string | |
| address | string | |
| debt | float | running total owed |
| status | string | ACTIVE / INACTIVE |

### imports
Replaces Purchase Orders — simple stock-in log, no supplier concept.
| field | type | note |
|---|---|---|
| id | ObjectId | |
| productId | ObjectId | ref products |
| quantity | float | added to product.stockQty on create |
| costPrice | float | may also update product.costPrice |
| userId | ObjectId | who recorded it |
| note | string | optional |
| date | datetime | |

### sales
| field | type | note |
|---|---|---|
| id | ObjectId | |
| code | string | unique, renamed from invoice_no |
| customerId | ObjectId | ref customers, optional |
| userId | ObjectId | ref users, cashier |
| saleDate | datetime | |
| total | float | renamed from grand_total |
| paid | float | renamed from paid_amount |
| debt | float | renamed from debt_amount |
| paymentStatus | string | PAID / UNPAID / PARTIAL |
| dueDate | date | optional |
| paymentMethod | string | CASH / TRANSFER / etc |

### saleItems
Separate collection (not embedded).
| field | type | note |
|---|---|---|
| id | ObjectId | |
| saleId | ObjectId | ref sales |
| productId | ObjectId | ref products |
| quantity | float | |
| costPrice | float | snapshot at sale time |
| unitPrice | float | snapshot at sale time |
| priceType | string | RETAIL / WHOLESALE, auto-decided by wholesaleMinQty |
| subtotal | float | quantity * unitPrice |

### debtPayments
| field | type | note |
|---|---|---|
| id | ObjectId | |
| saleId | ObjectId | ref sales |
| customerId | ObjectId | ref customers |
| userId | ObjectId | who received payment |
| amountPaid | float | |
| paymentDate | datetime | |
| note | string | optional |

### damagedProducts
| field | type | note |
|---|---|---|
| id | ObjectId | |
| productId | ObjectId | ref products |
| userId | ObjectId | |
| quantity | float | |
| costPrice | float | snapshot |
| reason | string | EXPIRED / BROKEN / LEAKED / OTHER |
| note | string | optional |
| date | datetime | |

### expenses
| field | type | note |
|---|---|---|
| id | ObjectId | |
| userId | ObjectId | |
| title | string | e.g. ຄ່ານໍ້າມັນ, ຄ່າຊື້ສິນຄ້າ |
| type | string | EXPENSE / INCOME |
| amount | float | |
| date | datetime | |

## Business rules

- Wholesale pricing: at sale time, if `quantity >= product.wholesaleMinQty` then `unitPrice = product.wholesalePrice` and `priceType = WHOLESALE`, else `unitPrice = product.retailPrice` and `priceType = RETAIL`. Decided automatically, not chosen manually by cashier.
- Restocking: no separate UI — done either via the `imports` collection (preferred, keeps a log) or by editing a product's `stockQty` directly.
- No stock movement audit log — `stockQty` is the single source of truth, updated directly by sales, imports, and damagedProducts.

## Reports (computed, not stored)

- **Cost/Profit**: from `saleItems` (unitPrice - costPrice) * quantity, minus `expenses` (EXPENSE), plus `expenses` (INCOME), minus `damagedProducts` (costPrice * quantity)
- **Sales total**: sum of `sales.total` over a date range
- **Unpaid invoices count**: count of `sales` where `paymentStatus != PAID`
- **Total debt**: sum of `sales.debt` (or sum of `customers.debt`)

Open question: reports computed real-time on query vs. cached/precomputed — decide before implementing report resolvers.
