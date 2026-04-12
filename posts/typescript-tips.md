# TypeScript の便利な型テクニック
date: 2024-03-20

TypeScript を書くうえで知っておくと便利な型のテクニックをまとめました。

## Mapped Types

オブジェクトの全プロパティを変換する型です。

```typescript
type Readonly<T> = {
  readonly [K in keyof T]: T[K];
};

type Partial<T> = {
  [K in keyof T]?: T[K];
};

// 使用例
interface User {
  name: string;
  age: number;
  email: string;
}

type ReadonlyUser = Readonly<User>;
type PartialUser = Partial<User>;
```

## Template Literal Types

文字列リテラル型を組み合わせた型です。

```typescript
type EventName = "click" | "focus" | "blur";
type Handler = `on${Capitalize<EventName>}`;
// => "onClick" | "onFocus" | "onBlur"

type CSSProperty = "margin" | "padding";
type Side = "top" | "right" | "bottom" | "left";
type CSSKey = `${CSSProperty}-${Side}`;
// => "margin-top" | "margin-right" | ...
```

## infer キーワード

条件型の中で型を推論します。

```typescript
type ReturnType<T> = T extends (...args: any[]) => infer R ? R : never;

type UnwrapPromise<T> = T extends Promise<infer U> ? U : T;

// 使用例
async function fetchUser(): Promise<User> { /* ... */ }

type FetchResult = UnwrapPromise<ReturnType<typeof fetchUser>>;
// => User
```

## Discriminated Union

型の絞り込みに使うパターンです。

```typescript
type Result<T> =
  | { success: true; data: T }
  | { success: false; error: string };

function handleResult<T>(result: Result<T>) {
  if (result.success) {
    console.log(result.data); // T
  } else {
    console.error(result.error); // string
  }
}
```

## まとめ

TypeScript の型システムは非常に強力です。うまく活用することで、
実行時エラーを大幅に減らせます。
