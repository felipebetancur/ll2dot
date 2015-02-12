; ModuleID = 'funcs.ll'
target datalayout = "e-m:e-i64:64-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

; Function Attrs: nounwind uwtable
define i32 @bar(i32 %x) #0 {
  br label %1

; <label>:1                                       ; preds = %3, %0
  %.0 = phi i32 [ %x, %0 ], [ %4, %3 ]
  %2 = icmp slt i32 %.0, 1000
  br i1 %2, label %3, label %5

; <label>:3                                       ; preds = %1
  %4 = mul nsw i32 %.0, 2
  br label %1

; <label>:5                                       ; preds = %1
  ret i32 %.0
}

; Function Attrs: nounwind uwtable
define i32 @main(i32 %argc, i8** %argv) #0 {
  %1 = icmp slt i32 undef, 3
  br i1 %1, label %2, label %4

; <label>:2                                       ; preds = %0
  %3 = call i32 @bar(i32 undef)
  br label %7

; <label>:4                                       ; preds = %0
  %5 = mul nsw i32 undef, 2
  %6 = call i32 @bar(i32 %5)
  br label %7

; <label>:7                                       ; preds = %4, %2
  %.0 = phi i32 [ %3, %2 ], [ %6, %4 ]
  ret i32 %.0
}

attributes #0 = { nounwind uwtable "less-precise-fpmad"="false" "no-frame-pointer-elim"="true" "no-frame-pointer-elim-non-leaf" "no-infs-fp-math"="false" "no-nans-fp-math"="false" "stack-protector-buffer-size"="8" "unsafe-fp-math"="false" "use-soft-float"="false" }

!llvm.ident = !{!0}

!0 = metadata !{metadata !"clang version 3.5.1 (tags/RELEASE_351/final)"}
