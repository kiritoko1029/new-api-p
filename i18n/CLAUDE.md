[Root](../CLAUDE.md) > **i18n**

# i18n -- Backend Internationalization

## Changelog

| Date | Action | Summary |
|------|--------|---------|
| 2026-03-25 | Created | Initial module scan and documentation. |

## Module Responsibility

Provides backend internationalization using `nicksnyder/go-i18n`. Used for system emails, error messages, and other server-side strings.

## Languages

| Language | File |
|----------|------|
| English | `locales/en.yaml` |
| Simplified Chinese | `locales/zh-CN.yaml` |
| Traditional Chinese | `locales/zh-TW.yaml` |

## Usage

- `i18n.Init()` -- Load translation files and initialize
- `i18n.SupportedLanguages()` -- List available languages
- `i18n.SetUserLangLoader(fn)` -- Set a function to load user language preference lazily
- Translations are accessed via the go-i18n `Localizer` API

## FAQ

**Q: How do I add a new backend language?**
A: Create a new YAML file in `i18n/locales/`, add the language code to the supported languages list, and register it in `i18n.Init()`.

## Related Files

- `i18n/locales/en.yaml` -- English translations
- `i18n/locales/zh-CN.yaml` -- Simplified Chinese translations
- `i18n/locales/zh-TW.yaml` -- Traditional Chinese translations
