# UAD 語言範例程式

這個目錄包含了多個 UAD 語言的小程式範例，展示不同的語言特性。

## 基礎範例

### hello_world.uad

最簡單的 Hello World 程式，展示基本的函數定義和輸出。

### simple_calc.uad

簡單的計算器，展示基本的算術運算和函數定義。

### basic_math.uad

基礎數學運算範例。

## 控制流

### if_example.uad

展示 `if/else` 表達式的使用，包括 `max`、`min`、`abs` 函數。

### while_loop.uad

展示 `while` 迴圈的使用，包括倒數計時和累加運算。

### factorial.uad

階乘計算，展示遞迴和迭代兩種實現方式。

### fibonacci.uad

費波那契數列計算，展示遞迴函數。

## 資料結構

### struct_example.uad

展示結構體（struct）的定義和使用，包括欄位存取和函數參數。

### match_example.uad

展示 `enum` 和 `match` 表達式的組合使用，實現 Result 類型模式。

### option_type.uad

實現 Option 類型，展示安全的除法和錯誤處理模式。

## 陣列和字串

### array_operations.uad

陣列操作範例，包括求和、找最大值等操作。

### string_operations.uad

字串操作範例，展示字串處理函數。

## 進階範例

### calculator.uad

計算器範例，展示結構體、函數和浮點數運算，包括計算兩點距離和中點。

## 語言特性展示

這些範例展示了 UAD 語言的主要特性：

- ✅ **函數定義**：`fn` 關鍵字定義函數
- ✅ **類型系統**：靜態類型，包括 `Int`、`Float`、`String`、`Bool`
- ✅ **結構體**：`struct` 定義複合類型
- ✅ **列舉**：`enum` 定義和類型（sum types）
- ✅ **模式匹配**：`match` 表達式進行模式匹配
- ✅ **控制流**：`if/else`、`while` 迴圈
- ✅ **陣列**：陣列字面量和索引存取
- ✅ **表達式導向**：函數返回表達式，`if` 和 `match` 都是表達式

## 執行方式

目前編譯器還在開發中，這些程式可以通過以下方式測試：

```bash
# 使用 uadc 編譯器（目前僅解析階段）
./bin/uadc -i examples/core/hello_world.uad -o output.uadir
```

## 注意事項

- 這些程式展示了 UAD 語言的語法，但某些內建函數（如 `println`、`sqrt`、`len` 等）的實現取決於標準庫
- 字串連接使用 `+` 運算符
- 陣列索引使用 `arr[i]` 語法
- enum 值直接使用 variant 名稱，不需要前綴（如 `Ok(42)` 而不是 `Result::Ok(42)`）
