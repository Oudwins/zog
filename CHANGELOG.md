# Changelog

## [0.21.5](https://github.com/Oudwins/zog/compare/v0.21.4...v0.21.5) (2025-07-23)


### Bug Fixes

* z.time.Format() not working ([638ff57](https://github.com/Oudwins/zog/commit/638ff574c03e4cbf1ab3d4f2e28d3300a8bfa2c1))

## [0.21.4](https://github.com/Oudwins/zog/compare/v0.21.3...v0.21.4) (2025-07-13)


### Features

* like primitive types ([#172](https://github.com/Oudwins/zog/issues/172)) ([126d79b](https://github.com/Oudwins/zog/commit/126d79b72c9a0ca9c736761effda7a45d7451222))

## [0.21.3](https://github.com/Oudwins/zog/compare/v0.21.2...v0.21.3) (2025-07-03)


### Features

* Add Azerbaijani language to i18n package ([#164](https://github.com/Oudwins/zog/issues/164)) ([6b75bad](https://github.com/Oudwins/zog/commit/6b75bad58c4e9901ee3f57f3ba9ba4f8a63f57e2))
* uint ([#171](https://github.com/Oudwins/zog/issues/171)) ([97ce447](https://github.com/Oudwins/zog/commit/97ce44749d5b5077c86693583aed819140925aaf))

## [0.21.2](https://github.com/Oudwins/zog/compare/v0.21.1...v0.21.2) (2025-06-15)


### Bug Fixes

* error message for number().oneOf() ([5de48a3](https://github.com/Oudwins/zog/commit/5de48a3a6537271d2ba046b5ec3ccf9c5f9db95a))
* optional in zhttp ([#167](https://github.com/Oudwins/zog/issues/167)) ([f1e621b](https://github.com/Oudwins/zog/commit/f1e621b65d07b5bd58de15faa2df9702c2e3b075))
* regex for email ([#166](https://github.com/Oudwins/zog/issues/166)) ([954de17](https://github.com/Oudwins/zog/commit/954de178191db19302bc379cd2050b62b7bbcfba))

## [0.21.1](https://github.com/Oudwins/zog/compare/v0.21.0...v0.21.1) (2025-05-24)


### Features

* add other not operators ([#155](https://github.com/Oudwins/zog/issues/155)) ([16a63a4](https://github.com/Oudwins/zog/commit/16a63a4c2e432c360bf30871859fd630160d4f62))

## [0.21.0](https://github.com/Oudwins/zog/compare/v0.20.0...v0.21.0) (2025-05-03)


### ⚠ BREAKING CHANGES

* transforms and tests are now typesafe for the schema you are using. No more typecasting in Zog! ([#149](https://github.com/Oudwins/zog/issues/149))

### Features

* provide better panic msgs and docs ([#153](https://github.com/Oudwins/zog/issues/153)) ([f605689](https://github.com/Oudwins/zog/commit/f605689511eabed2343e333379607bb0aea93937))
* transforms and tests are now typesafe for the schema you are using. No more typecasting in Zog! ([#149](https://github.com/Oudwins/zog/issues/149)) ([2da0c0c](https://github.com/Oudwins/zog/commit/2da0c0c174a4976c6ab2ac33a05477d0c8454bcb))

## [0.20.0](https://github.com/Oudwins/zog/compare/v0.19.2...v0.20.0) (2025-04-26)


### ⚠ BREAKING CHANGES

* Deprecated z.Schema in favor of z.Shape. No plans to remove z.Schema for now
* Allow empty strings during parsing to not trigger required. Instead do .Min(1). ([#148](https://github.com/Oudwins/zog/issues/148))
* Transforms are now run sequentially with tests in the order they were defined in. schema.PostTransform has been deprecated and until it is removed it will work just like schema.Transform. Therefore schema.Min(1).Trim().Min(1) will work and first check string length > 1 then trim then check again for string len > 1
* implemented preprocess, removed preTransforms

### Features

* Deprecated z.Schema in favor of z.Shape. No plans to remove z.Schema for now ([544f03d](https://github.com/Oudwins/zog/commit/544f03d092b1ff3e46a52f9cb246a0211d44964a))
* implemented preprocess, removed preTransforms ([9dde9be](https://github.com/Oudwins/zog/commit/9dde9be47a11c5a584f85afd72ddb167bcc71325))
* ordered transforms ([#147](https://github.com/Oudwins/zog/issues/147)) ([1db8140](https://github.com/Oudwins/zog/commit/1db81403e1a77bf74a024be59e109882a9d0c91c))


### Bug Fixes

* Allow empty strings during parsing to not trigger required. Instead do .Min(1). ([#148](https://github.com/Oudwins/zog/issues/148)) ([da1b680](https://github.com/Oudwins/zog/commit/da1b6802bc13acd604b7d99f2071f35ad891084d))

## [0.19.2](https://github.com/Oudwins/zog/compare/v0.19.1...v0.19.2) (2025-04-18)


### Features

* add not operator ([5049e4e](https://github.com/Oudwins/zog/commit/5049e4ec1555458e7596bb3ebda485397516e9f5))
* add not operator ([5049e4e](https://github.com/Oudwins/zog/commit/5049e4ec1555458e7596bb3ebda485397516e9f5))

## [0.19.1](https://github.com/Oudwins/zog/compare/v0.19.0...v0.19.1) (2025-04-13)


### Features

* customFunc for easy custom schemas. Usage is z.CustomFunc[T any](func (valPtr *T, ctx z.Ctx) bool {}, ...TestOptions) ([#141](https://github.com/Oudwins/zog/issues/141)) ([4c2b42e](https://github.com/Oudwins/zog/commit/4c2b42e2bfa243347965dbb8f6d9d01c203e1699))


### Miscellaneous Chores

* release 0.19.1 ([0485a6b](https://github.com/Oudwins/zog/commit/0485a6b4a4a4eae98db24000239e738730726fec))

## [0.19.0](https://github.com/Oudwins/zog/compare/v0.18.4...v0.19.0) (2025-04-03)


### ⚠ BREAKING CHANGES

* Test.ValidateFunc's name has been changed to Test.Func

### Features

* implemented super refine like API. Create complex custom tests easily ([#138](https://github.com/Oudwins/zog/issues/138)) ([e44e593](https://github.com/Oudwins/zog/commit/e44e593034be7f91f61c63193e5848bc0896a60b))


### Miscellaneous Chores

* release 0.19.0 ([2a25dc4](https://github.com/Oudwins/zog/commit/2a25dc45792cb4f0ea35cb999e6b3c6ea19f0c88))

## [0.18.4](https://github.com/Oudwins/zog/compare/v0.18.2...v0.18.4) (2025-03-16)


### Features

* support for custom strings, numbers and booleans ([#131](https://github.com/Oudwins/zog/issues/131)) ([29cb24d](https://github.com/Oudwins/zog/commit/29cb24db55ba520a65e57c4d2e6ca729fffddbc2))


### Miscellaneous Chores

* release 0.18.4 ([a3526ae](https://github.com/Oudwins/zog/commit/a3526aea1eaace1738bb41d0f82b5c6b9014e9c5))

## [0.18.2](https://github.com/Oudwins/zog/compare/v0.18.1...v0.18.2) (2025-03-08)


### ⚠ BREAKING CHANGES

* zog tag changed to a catch all instead of superseding the other tags

### Bug Fixes

* zog tag changed to a catch all instead of superseding the other tags ([869268e](https://github.com/Oudwins/zog/commit/869268ecf8aaf981fe0ac23e7de521bc22222e2a))


### Miscellaneous Chores

* release 0.18.2 ([6517301](https://github.com/Oudwins/zog/commit/6517301b14f44b40c73e5193e846fdb449005512))

## [0.18.1](https://github.com/Oudwins/zog/compare/v0.18.0...v0.18.1) (2025-03-08)


### Features

* allow zog to use multiple struct tags zog, json, form, query, env ([5426bb6](https://github.com/Oudwins/zog/commit/5426bb6eb5c93f314dc0c5d805d5e7cd51b2cf0c))


### Miscellaneous Chores

* release 0.18.1 ([62fcb88](https://github.com/Oudwins/zog/commit/62fcb88e420204837ee23902c6761037d5fafc3d))

## [0.18.0](https://github.com/Oudwins/zog/compare/v0.17.2...v0.18.0) (2025-03-08)


### ⚠ BREAKING CHANGES

* DELETE Method on zhttp also supports multiplexing based on content type header
* Removed deprecated z.ParseCtx interface. Use z.Ctx instead
* Zog issue interface is converted into a struct ([#119](https://github.com/Oudwins/zog/issues/119))

### Bug Fixes

* DELETE Method on zhttp also supports multiplexing based on content type header ([b8ebdc7](https://github.com/Oudwins/zog/commit/b8ebdc71856a8de3fc59cac789f592a4900c7333))
* export pre & post transforms from main zog package ([57703a4](https://github.com/Oudwins/zog/commit/57703a4e8c4b6ab3c86b80e6cfa4f81d2e5ddda5))


### Miscellaneous Chores

* release 0.18.0 ([df12fb1](https://github.com/Oudwins/zog/commit/df12fb1144722a62b8bea492b3101f674532ebbd))


### Code Refactoring

* Removed deprecated z.ParseCtx interface. Use z.Ctx instead ([0bb4cfd](https://github.com/Oudwins/zog/commit/0bb4cfde4df4f5b8a72e2772e398f17331d2a235))
* Zog issue interface is converted into a struct ([#119](https://github.com/Oudwins/zog/issues/119)) ([a68b393](https://github.com/Oudwins/zog/commit/a68b393d7b078274c26bad459340a1d7276deef1))

## [0.17.2](https://github.com/Oudwins/zog/compare/v0.17.1...v0.17.2) (2025-03-02)


### Performance Improvements

* removed allocation for struct schemas ([4cd379c](https://github.com/Oudwins/zog/commit/4cd379c114877200105193a00fd76d08197b4236))

## [0.17.1](https://github.com/Oudwins/zog/compare/v0.17.0...v0.17.1) (2025-02-28)


### Performance Improvements

* removed one allocation per schema ([14f3511](https://github.com/Oudwins/zog/commit/14f3511bb2a18750d429133c77504c05e09d2777))

## [0.17.0](https://github.com/Oudwins/zog/compare/v0.16.6...v0.17.0) (2025-02-22)


### Features

* new schemas, float32, int64, int32 ([#108](https://github.com/Oudwins/zog/issues/108)) ([0e57a5b](https://github.com/Oudwins/zog/commit/0e57a5b7ef591ebf251176c6532490abcf3b78e1))
* number schema now supports any number ([#110](https://github.com/Oudwins/zog/issues/110)) ([6af9605](https://github.com/Oudwins/zog/commit/6af96054200e26bb18077dc3bd384be813296d3c))

## [0.16.6](https://github.com/Oudwins/zog/compare/v0.16.5...v0.16.6) (2025-02-21)


### Bug Fixes

* delete method also only allows for parsing of query params just like GET and HEAD with zhttp ([6bc7636](https://github.com/Oudwins/zog/commit/6bc763695061683acb51fbf81c1ae965e2d43546))
* zhttp GET request with json or form content type still fetches from params ([867cd06](https://github.com/Oudwins/zog/commit/867cd063414e622e0561a2e302a34999e621adab))


### Performance Improvements

* slice pathbuilder for fewer allocations ([a20f5e2](https://github.com/Oudwins/zog/commit/a20f5e2678007b20832f4d38e1e3683ca49d95e9))
* slice pathbuilder for fewer allocations ([a20f5e2](https://github.com/Oudwins/zog/commit/a20f5e2678007b20832f4d38e1e3683ca49d95e9))

## [0.16.5](https://github.com/Oudwins/zog/compare/v0.16.4...v0.16.5) (2025-02-20)


### Performance Improvements

* removed string sync pool since it has worse performance ([#103](https://github.com/Oudwins/zog/issues/103)) ([602c146](https://github.com/Oudwins/zog/commit/602c14657bc683ffc920f830d7b6092b29bd2efd))

## [0.16.4](https://github.com/Oudwins/zog/compare/v0.16.3...v0.16.4) (2025-02-20)


### Bug Fixes

* panic on single letter query param zhttp ([#101](https://github.com/Oudwins/zog/issues/101)) ([cd8d172](https://github.com/Oudwins/zog/commit/cd8d172d30723dafd18167d598e8c6f70417e3e9))
* zhttp supports more complex content type strings ([#99](https://github.com/Oudwins/zog/issues/99)) ([9460ea2](https://github.com/Oudwins/zog/commit/9460ea285d91dc50a973d7a4d83bcf754559c0b3))


### Performance Improvements

* string builder for default langmap replaceble placeholders ([#102](https://github.com/Oudwins/zog/issues/102)) ([8f3c881](https://github.com/Oudwins/zog/commit/8f3c881530995960d98132e1f96a2d582f742b58))

## [0.16.3](https://github.com/Oudwins/zog/compare/v0.16.2...v0.16.3) (2025-02-20)


### Performance Improvements

* internal tests for primitive values now use pointers. User tests are unaffected ([#97](https://github.com/Oudwins/zog/issues/97)) ([6d3a234](https://github.com/Oudwins/zog/commit/6d3a2345b49461b6b089e4f91366760940afec05))

## [0.16.2](https://github.com/Oudwins/zog/compare/v0.16.1...v0.16.2) (2025-02-19)


### Performance Improvements

* issues syncpool and issues.Collect function to collect issues into sync pool ([#93](https://github.com/Oudwins/zog/issues/93)) ([91d6b3e](https://github.com/Oudwins/zog/commit/91d6b3e49bf11af4cfcbd97c0a02be5da478504f))

## [0.16.1](https://github.com/Oudwins/zog/compare/v0.16.0...v0.16.1) (2025-02-19)


### Features

* Params Option ([#90](https://github.com/Oudwins/zog/issues/90)) ([05823ab](https://github.com/Oudwins/zog/commit/05823abb1658ae70921b294efb2d498afb913278))


### Performance Improvements

* implemented syncpool for internal structures ([#91](https://github.com/Oudwins/zog/issues/91)) ([f1717f6](https://github.com/Oudwins/zog/commit/f1717f62ad5a9ad4f27c9e8de4f194460dd3f3a5))


### Miscellaneous Chores

* release 0.16.1 ([6d416fa](https://github.com/Oudwins/zog/commit/6d416faac3d86c1ddd62743cfc6323fe915b51a1))

## [0.16.0](https://github.com/Oudwins/zog/compare/v0.15.1...v0.16.0) (2025-02-11)


### ⚠ BREAKING CHANGES

* structs can no longer be required or optional. Define this in the fields instead. If you need to model a struct that might exist use a pointer to a struct. This should not affect most users as now it works how everyone intuitively thought it worked. ([#88](https://github.com/Oudwins/zog/issues/88))
* renamed ZogError to ZogIssue to be more aligned with Zod. Deprecated a bunch of APIs for naming consistency. conf.ErrorFormatter removed in favor of conf.IssueFormatter ([#86](https://github.com/Oudwins/zog/issues/86))

### Features

* new test option to set the issue path for issues generated from that test ([#87](https://github.com/Oudwins/zog/issues/87)) ([47d24a1](https://github.com/Oudwins/zog/commit/47d24a115fd94198447fed8df04690bf019f305c))
* testFunc method on schemas for easier custom tests ([#82](https://github.com/Oudwins/zog/issues/82)) ([52a90eb](https://github.com/Oudwins/zog/commit/52a90eb197b5380499319c1e29cf61ae1665e3e1))


### Bug Fixes

* structs can no longer be required or optional. Define this in the fields instead. If you need to model a struct that might exist use a pointer to a struct. This should not affect most users as now it works how everyone intuitively thought it worked. ([#88](https://github.com/Oudwins/zog/issues/88)) ([9681ebb](https://github.com/Oudwins/zog/commit/9681ebb691a2cfa188f0fc539024cec2cacfcaa3))


### Miscellaneous Chores

* release 0.16.0 ([3a222d0](https://github.com/Oudwins/zog/commit/3a222d06fc3ff463aa958165dd23b0cf612ec9a3))


### Code Refactoring

* renamed ZogError to ZogIssue to be more aligned with Zod. Deprecated a bunch of APIs for naming consistency. conf.ErrorFormatter removed in favor of conf.IssueFormatter ([#86](https://github.com/Oudwins/zog/issues/86)) ([49f01e3](https://github.com/Oudwins/zog/commit/49f01e3b6d522da9edb44bb076f4bd9baa7957f4))

## [0.15.1](https://github.com/Oudwins/zog/compare/v0.15.0...v0.15.1) (2025-02-09)


### Bug Fixes

* ZogErr Value() doesn't return underlying value ([#77](https://github.com/Oudwins/zog/issues/77)) ([1f7d2ff](https://github.com/Oudwins/zog/commit/1f7d2ff58b9d5a02b6bdd5a70c709e8dc40eadc5))

## [0.15.0](https://github.com/Oudwins/zog/compare/v0.14.1...v0.15.0) (2025-02-03)


### Features

* New Validate method to validate existing structures ([#68](https://github.com/Oudwins/zog/issues/68)) ([8297022](https://github.com/Oudwins/zog/commit/82970228b25b3eaa619dc24e0237d754330d6b28))


### Miscellaneous Chores

* release 0.15.0 ([17cc76d](https://github.com/Oudwins/zog/commit/17cc76d1229e96501c25190f6b9a4ffadf3b2ef5))

## [0.14.1](https://github.com/Oudwins/zog/compare/v0.14.0...v0.14.1) (2025-01-13)


### Bug Fixes

* Schemas are now public. New complex & primitive schema interfaces ([#64](https://github.com/Oudwins/zog/issues/64)) ([9e659e5](https://github.com/Oudwins/zog/commit/9e659e50b508418c1faaeccf7cc57b7c3b1cbb98))

## [0.14.0](https://github.com/Oudwins/zog/compare/v0.13.0...v0.14.0) (2025-01-02)


### Features

* zjson package ([#58](https://github.com/Oudwins/zog/issues/58)) ([2009561](https://github.com/Oudwins/zog/commit/20095613cbe017772bb92e774cffcb4f8cda390b))

## [0.13.0](https://github.com/Oudwins/zog/compare/v0.12.1...v0.13.0) (2024-11-12)


### Features

* Implemented Struct().Pick(), Struct().Omit() and Struct().Extend() ([#53](https://github.com/Oudwins/zog/issues/53)) ([8adc803](https://github.com/Oudwins/zog/commit/8adc80356997c34725d11ee03905bd843f742067))

## [0.12.1](https://github.com/Oudwins/zog/compare/v0.12.0...v0.12.1) (2024-11-11)


### Features

* support for parsing into pointers. Now you may have pointers in the destination ([#42](https://github.com/Oudwins/zog/issues/42)) ([fd6bbbf](https://github.com/Oudwins/zog/commit/fd6bbbf5f5afb7fbafdba066e43c6bf059f3d6b6))


### Miscellaneous Chores

* release 0.12.1 ([505dcdb](https://github.com/Oudwins/zog/commit/505dcdb181269042b27fad80ad350dfefa49961e))

## [0.12.0](https://github.com/Oudwins/zog/compare/v0.11.0...v0.12.0) (2024-11-09)


### Features

* implement z.String().Trim() as a built in PreTransform that trims the input data if it is a string ([#51](https://github.com/Oudwins/zog/issues/51)) ([1d65859](https://github.com/Oudwins/zog/commit/1d65859a9c906ad5905220d06b2b0e3c3d9c628d))
* schema custom coercer support via the z.WithCoercer function and custom time formats via z.Time.Format() fuction ([#48](https://github.com/Oudwins/zog/issues/48)) ([1472669](https://github.com/Oudwins/zog/commit/1472669a66b2928a18a923e794faa373821961cd))
* time coercer now support for unix timestamps in ms  ([#47](https://github.com/Oudwins/zog/issues/47)) ([4c5b4bd](https://github.com/Oudwins/zog/commit/4c5b4bd56ac762325875c4adf0dbb9021e72fe00))
* zhttp package now supports providing your own custom parsers ([#50](https://github.com/Oudwins/zog/issues/50)) ([e8a111f](https://github.com/Oudwins/zog/commit/e8a111fdad679b19a9b01968eff846de777ad24a))


### Bug Fixes

* required check not working with zero values from other types ([#44](https://github.com/Oudwins/zog/issues/44)) ([1abc8e8](https://github.com/Oudwins/zog/commit/1abc8e853f631a2bd34f86279140344228740371))

## [0.11.0](https://github.com/Oudwins/zog/compare/v0.10.0...v0.11.0) (2024-11-01)


### ⚠ BREAKING CHANGES

* zhttp.NewRequestDataProvider() which was deprecated is now removed. Please use zhttp.Request() instead
* ZogError.Error() no longer proxies to the wrapped error. Now it returns a string representation of the ZogError. You can still access Wrapped error through Unwrap()

### Features

* Add string regex and uuid validators  ([#40](https://github.com/Oudwins/zog/issues/40)) ([8853d76](https://github.com/Oudwins/zog/commit/8853d76010adf109a9e912c306ee4811d3b62155))
* error printing ([#43](https://github.com/Oudwins/zog/issues/43)) ([5446312](https://github.com/Oudwins/zog/commit/54463121bb9425f20e2522b790629ac1894f5268))


### Bug Fixes

* consider "    " to be a zero value ([d4856e9](https://github.com/Oudwins/zog/commit/d4856e9bb0c973e7df3bb653111bc14a361b69df))
* zhttp handles input json being null ([9e8b8d3](https://github.com/Oudwins/zog/commit/9e8b8d3696d1c8fb505fcaf590b43c9f4150aa1d))


### Miscellaneous Chores

* release ([6bfc8cc](https://github.com/Oudwins/zog/commit/6bfc8cc22ea9cc76bd26e72dfa64dc5e639bd3a7))


### Code Refactoring

* removed zhttp new data provider ([86d4e6f](https://github.com/Oudwins/zog/commit/86d4e6f9c10654af11f1a7e4d0d1ba7ccddd5245))

## [0.10.0](https://github.com/Oudwins/zog/compare/v0.9.1...v0.10.0) (2024-10-07)


### Features

* add test options to time methods ([42db318](https://github.com/Oudwins/zog/commit/42db318598f4673b8369925d82ebc57b39dba0c5))


### Bug Fixes

* boolean false parsing behavior ([#35](https://github.com/Oudwins/zog/issues/35)) ([8670b64](https://github.com/Oudwins/zog/commit/8670b641d53724bed570a89c2c60da8dd1645b13))
* struct merge panic on merging schemas with transforms ([90ccc88](https://github.com/Oudwins/zog/commit/90ccc885cc4579dcf80bcce3743e0d4992b22d14))

## [0.9.1](https://github.com/Oudwins/zog/compare/v0.9.0...v0.9.1) (2024-09-26)


### Bug Fixes

* zhttp request method ([f69af8b](https://github.com/Oudwins/zog/commit/f69af8b09107236878388abda005966e8b16726c))

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
