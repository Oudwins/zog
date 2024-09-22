# Changelog

## [0.9.0](https://github.com/Oudwins/zog/compare/v0.8.0...v0.9.0) (2024-09-22)


### Features

* i18n package with spanish & english translations ([#28](https://github.com/Oudwins/zog/issues/28)) ([1120fd6](https://github.com/Oudwins/zog/commit/1120fd68759696f1ef9113757af2b5448fa10448))
* improved zhttp library ([#32](https://github.com/Oudwins/zog/issues/32)) ([891bb6c](https://github.com/Oudwins/zog/commit/891bb6c4c50896094c6c370832486b0862d1172a))
* trim space for env variables ([166d881](https://github.com/Oudwins/zog/commit/166d8812e429ea8648a9272006d29290fb227733))

## [0.8.0](https://github.com/Oudwins/zog/compare/v0.7.0...v0.8.0) (2024-09-16)


### Features

* added transforms for slices ([#22](https://github.com/Oudwins/zog/issues/22)) ([71c01cb](https://github.com/Oudwins/zog/commit/71c01cb0b97741f35b9ae1bdce8f4ac966881b41))
* json data provider ([2d10c92](https://github.com/Oudwins/zog/commit/2d10c92eeb4d273e11d2bd4ede8e1b1741e897c4))
* panic on invalid struct schema ([#25](https://github.com/Oudwins/zog/issues/25)) ([8d5b493](https://github.com/Oudwins/zog/commit/8d5b49329102d6db6e7fa949cb3228e7fe9ee72f))


### Bug Fixes

* required custom z.Message ([#24](https://github.com/Oudwins/zog/issues/24)) ([49198a0](https://github.com/Oudwins/zog/commit/49198a0b8d9678f76fd87b4328122a6392a902be))
* structs now handle both uppercase and lowercase first letters ([4fbc9c3](https://github.com/Oudwins/zog/commit/4fbc9c305cd3d91b224753d3ba76e7d7ba250b6b))

## [0.7.0](https://github.com/Oudwins/zog/compare/v0.6.2...v0.7.0) (2024-09-09)


### ⚠ BREAKING CHANGES

* custom tests now require that you pass a test struct or use the TestFunc() helper
* order of schema.Test() params has changed from (errorCode, z.Message(), func) to (errCode, func, [optionalTestOptions])
* All z.Errors functions have changed. I still don't recommend you use them since they might still change in the future

### Features

* better errors ([fe78a8d](https://github.com/Oudwins/zog/commit/fe78a8d072abf23f9c7d60d2b8560d2384dd899f))
* move coercers to default variable to make it easier to replace the coercers struct without losing access to the default coercers ([2387330](https://github.com/Oudwins/zog/commit/2387330e60306ed7767a90be155d7479df12e21c))
* new & improved API for custom tests ([9acfc37](https://github.com/Oudwins/zog/commit/9acfc378a530f651b2bcf2e7c4f344f4cbc2f8d2))


### Bug Fixes

* bool coercer ([#14](https://github.com/Oudwins/zog/issues/14)) ([01f8c17](https://github.com/Oudwins/zog/commit/01f8c17b38050604112f688006b259a60df6a58a))
* minor fix to order of operations when required is set ([cff0fc3](https://github.com/Oudwins/zog/commit/cff0fc3a87bbc2c574601fe81daddafb5f48a279))
* Time().EQ() was broken due to typo ([9310e1a](https://github.com/Oudwins/zog/commit/9310e1a5dab72fac18440921a84b3fa68d65e9b0))


### Miscellaneous Chores

* release 0.7.0 ([0e0eb47](https://github.com/Oudwins/zog/commit/0e0eb47d8094f7f84f9581630534d2e26838bef9))


### Code Refactoring

* custom test method is now more in line with the rest. ([d163f36](https://github.com/Oudwins/zog/commit/d163f369ba6310cf849d1271a214ad95082bd641))

## [0.6.2](https://github.com/Oudwins/zog/compare/v0.6.1...v0.6.2) (2024-08-16)


### ⚠ BREAKING CHANGES

* slice errMap will now access validation errors for the first element through `[0]` key rather than `0` key

### Features

* slices now support structs ([#10](https://github.com/Oudwins/zog/issues/10)) ([52009ec](https://github.com/Oudwins/zog/commit/52009ec080aeff39c4904d7550d43d7fc84e33cd))


### Miscellaneous Chores

* release 0.6.2 ([b4c90c4](https://github.com/Oudwins/zog/commit/b4c90c4f98f91dc0602932adf364263319af9358))

## [0.6.1](https://github.com/Oudwins/zog/compare/v0.6.0...v0.6.1) (2024-08-16)


### Bug Fixes

* more realistic min go version ([#6](https://github.com/Oudwins/zog/issues/6)) ([658f060](https://github.com/Oudwins/zog/commit/658f060a66189ec0f0172cef507ab5e628442ed4))

## [0.6.0](https://github.com/Oudwins/zog/compare/v0.5.0...v0.6.0) (2024-08-16)


### Features

* quality of life improvements for working with errors ([#3](https://github.com/Oudwins/zog/issues/3)) ([1f3c3d0](https://github.com/Oudwins/zog/commit/1f3c3d003934fc4c9af81e9adabc20bce4e0fc8a))

## 0.5.0 (2024-08-12)


### Features

* added global functions to time validation ([d4abdca](https://github.com/Oudwins/zog/commit/d4abdcad414febb1372cfe61756f331041a6fd63))
* v0.5 release! ([#1](https://github.com/Oudwins/zog/issues/1)) ([7ac74c7](https://github.com/Oudwins/zog/commit/7ac74c72f9b5f59b87c561ce50d377b765c9a082))


### Bug Fixes

* better huristic for zhttp ([e311d91](https://github.com/Oudwins/zog/commit/e311d91610b72a65fdb3516fbc90f309c796b353))
* optional bug with slices ([107f4d6](https://github.com/Oudwins/zog/commit/107f4d694426ac0be3e4bc94cfcfa8f4d79aabea))


### Miscellaneous Chores

* release 0.5.0 ([6da4503](https://github.com/Oudwins/zog/commit/6da4503889dff7d32b7cc99344ba28e2ebd0da1c))
