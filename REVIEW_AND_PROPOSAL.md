# Review and Proposal: BotPlatform Refactoring

## 1. Review of Current State

The `BotPlatform` repository is currently tightly coupled with the `ScriptureBot` application. This coupling prevents `BotPlatform` from being a truly "democratized" and generic platform for other chatbots.

### Key Issues Identified:

1.  **Data Structure Coupling (`pkg/def/class.go`)**:
    *   **`UserData`**: Contains `datastore:""` tags. These are specific to Google Cloud Datastore and the schema used by `ScriptureBot`. A generic platform should be storage-agnostic.
    *   **`SessionData`**: Contains `ResourcePath string`. This is a `ScriptureBot`-specific configuration used to locate local resources. Generic session data should not enforce specific configuration fields.

2.  **Platform Implementation (`pkg/platform/telegram.go`)**:
    *   The `Translate` method populates `env.User` directly into the struct with `datastore` tags. While this works for `ScriptureBot`, it forces any other user of this library to use the same user structure or ignore the tags.

3.  **ScriptureBot Usage**:
    *   `ScriptureBot` relies on `BotPlatform`'s `UserData` for its database operations (`utils.RegisterUser`, `utils.PushUser`).
    *   `ScriptureBot` uses `SessionData.ResourcePath` to pass configuration to its command handlers (`app.ProcessCommand`).

## 2. Refactoring Proposal for BotPlatform

The goal is to remove all `ScriptureBot`-specific artifacts from `BotPlatform` while providing extension points so `ScriptureBot` (and other bots) can still function effectively.

### Proposed Changes:

1.  **Clean `UserData`**:
    *   Remove all `datastore` tags from the `UserData` struct in `pkg/def/class.go`.
    *   Ensure `UserData` only contains fields relevant to the chat platform identity (Id, Username, Firstname, Lastname, Type).

2.  **Generalize `SessionData`**:
    *   Remove `ResourcePath` from `SessionData`.
    *   Add a generic `Props map[string]interface{}` (or similar context mechanism) to `SessionData`. This allows applications to attach arbitrary data (like `ResourcePath` or other context) to the session as it flows through the pipeline.

3.  **Interface Stability**:
    *   Keep the `Platform` interface (`Translate`, `Post`) as is, as it is already reasonably generic.

## 3. Adaptation Plan for ScriptureBot

Since `BotPlatform` will be modifying its public API (struct definitions), `ScriptureBot` will need to be updated to consume the new version.

### Required Changes in ScriptureBot:

1.  **Define Local User Model**:
    *   Create a `User` struct in `ScriptureBot` (e.g., in `pkg/models/user.go`) that mirrors the fields needed but includes the `datastore` tags.
    *   Example:
        ```go
        type User struct {
            Firstname string `datastore:""`
            // ... other fields
        }
        ```

2.  **Map Data**:
    *   In `TelegramHandler` (and other entry points), after calling `bot.Translate` which returns `platform.SessionData`, map the `platform.UserData` to `ScriptureBot.User`.
    *   Use this local `User` struct for all database operations (`utils.RegisterUser`, `utils.PushUser`).

3.  **Handle ResourcePath**:
    *   Instead of `env.ResourcePath`, populates the new `Props` map in `SessionData`:
        ```go
        env.Props["ResourcePath"] = "/go/bin/"
        ```
    *   Update `app.ProcessCommand` and other consumers to read `ResourcePath` from `env.Props`.

4.  **Update Dependencies**:
    *   Update `go.mod` to point to the new version of `BotPlatform`.

## 4. Conclusion

This refactoring separates the concerns of "Platform Integration" (Telegram, Discord) from "Business Logic & Storage" (ScriptureBot). `BotPlatform` becomes a lightweight translation layer, and `ScriptureBot` retains full control over its data persistence and application context.
