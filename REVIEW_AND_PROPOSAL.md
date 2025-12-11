# Review and Proposal: BotPlatform Refactoring

## 1. Review of Current State

The `BotPlatform` repository is currently tightly coupled with the `ScriptureBot` application. This coupling prevents `BotPlatform` from being a truly "democratized" and generic platform for other chatbots.

### Key Issues Identified:

1.  **Data Structure Coupling (`pkg/def/class.go`)**:
    *   **`UserData`**: Contains `datastore:""` tags. These are specific to Google Cloud Datastore and the schema used by `ScriptureBot`. A generic platform should be storage-agnostic.
    *   **`UserData`**: Contains `Action` and `Config` fields. These are application-level state tracking fields specific to ScriptureBot's state machine logic, not properties of a Platform User.
    *   **`SessionData`**: Contains `ResourcePath string`. This is a `ScriptureBot`-specific configuration used to locate local resources. Generic session data should not enforce specific configuration fields.
    *   **UI Constraints**: The generic `ResponseOptions` struct forces a 1-column layout (via hardcoded constants in the Telegram implementation), limiting flexibility for other bots.

2.  **Platform Implementation (`pkg/platform/telegram.go`)**:
    *   The `Translate` method populates `env.User` directly into the struct. While functional, it needs to ensure generic extensibility points (like `Props`) are initialized.

3.  **ScriptureBot Usage**:
    *   `ScriptureBot` relies on `BotPlatform`'s `UserData` for its database operations (`utils.RegisterUser`, `utils.PushUser`) and state tracking (`Action`).
    *   `ScriptureBot` uses `SessionData.ResourcePath` to pass configuration.

## 2. Refactoring Proposal for BotPlatform

The goal is to remove all `ScriptureBot`-specific artifacts from `BotPlatform` while providing extension points so `ScriptureBot` (and other bots) can still function effectively.

### Proposed Changes:

1.  **Clean `UserData`**:
    *   Remove all `datastore` tags from the `UserData` struct.
    *   Remove `Action` and `Config` fields. `UserData` should only contain fields relevant to the chat platform identity (Id, Username, Firstname, Lastname, Type).

2.  **Generalize `SessionData`**:
    *   Remove `ResourcePath` from `SessionData`.
    *   Add a generic `Props map[string]interface{}`. This allows applications to attach arbitrary data (like `ResourcePath` or other context) to the session.
    *   **Crucial Implementation Detail**: Platform implementations (e.g., `Translate` in `telegram.go`) *must* initialize this map (`make(map[string]interface{})`) to prevent runtime panics for consumers.

3.  **Enhance UI Flexibility**:
    *   Add `ColWidth int` to `ResponseOptions`.
    *   Update platform logic to use this value for button layout, defaulting to the standard (1 column) if not set.

## 3. Adaptation Plan for ScriptureBot

Since `BotPlatform` will be modifying its public API, `ScriptureBot` must be updated.

### Required Changes in ScriptureBot:

1.  **Define Local User Model**:
    *   Create a `User` struct in `ScriptureBot` (e.g., in `pkg/models/user.go`) that includes:
        *   The basic fields (Firstname, etc.)
        *   The `datastore` tags.
        *   **The State Fields**: `Action` and `Config`.
    *   Example:
        ```go
        type User struct {
            Firstname string `datastore:""`
            Action    string `datastore:""`
            Config    string `datastore:""`
            // ... other fields
        }
        ```

2.  **Map Data**:
    *   In `TelegramHandler`, map `platform.UserData` (identity) to `ScriptureBot.User` (identity + state).
    *   Load `Action` and `Config` from the database (via `utils.RegisterUser`), not from the platform session.

3.  **Handle ResourcePath**:
    *   Populate `env.Props["ResourcePath"]` in the handler and read it from there in command processors.

## 4. Migration Impact Analysis

### Will this affect existing users?
**No, the data for existing users will remain intact.**

*   **Data Compatibility**: The removal of fields (`Action`, `Config`) from the *library struct* does not delete columns in the *database*. Since `ScriptureBot` will define a local struct that *includes* these fields before writing back to the DB, the data is preserved.
*   **Datastore Tags**: Removing tags is safe as the Go Datastore client defaults to field names, which matches the previous behavior.

### Do we need a migration task?
**Yes, a *Code Migration* task is required.**

`ScriptureBot` **will fail to compile** or **lose state functionality** without code changes because `UserData` will no longer have `Action` or `Config`.

*   **Task**: Implement the "Define Local User Model" step. This is critical to preserve the bot's ability to remember user state (e.g., "waiting for search term").

## 5. Conclusion

This refactoring strictly separates "Platform Identity" from "Application State" and "Storage". `BotPlatform` handles the delivery of messages, while `ScriptureBot` owns the user's state and data persistence.
