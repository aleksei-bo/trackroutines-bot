# TrackRoutines Bot for Telegram

The Telegram Bot is a simple bot written in Go using the `github.com/go-telegram-bot-api/telegram-bot-api/v5` library. It provides a way for users to track and manage their tasks through a Telegram chat interface. The bot also utilizes the `github.com/jmoiron/sqlx` and `github.com/mattn/go-sqlite3` packages for data storage and retrieval.

This is the initial version of the bot, so there may be some bugs or user experience issues. The bot stores data in a local SQLite database within the same folder as the bot's executable.

## About

The goal of this bot is to make completion of repetative tasks more interesting. Doesn't matter if it is a household routines (do dishes, cook, clean) or office routines, which repeats every day and are incredibly boring. With this bot you can create a list of tasks and assign values to each of them. Each completion of the task creates an action (due to the routine nature of the tasks, they are treated as repetative), and all values (points) of your actions are counted. It is easy to see who did the most stuff and who was giving some slack all the time.

## Features

- **Authentication**: You can set a custom authentication passphrase (`BOT_AUTH`) to prevent unauthorized access to the bot.
- **Task Management**: Users can add, view, and mark tasks as completed or uncompleted.
- **Task Persistence**: Tasks are stored in a local SQLite database, allowing users to access their task history across sessions.
- **Score**: Set a score for each task.

## Installation and Setup

To use the Telegram Task Count Bot, follow these steps:

1. **Create a Telegram Bot**: Use the BotFather account in Telegram (@BotFather) to create a new bot and obtain the bot token. Select `/newbot` from the list of commands and follow the instructions.
2. **Set the Bot Token**: Export the bot token as an environment variable named `BOT_TOKEN`.

   ```bash
   export BOT_TOKEN=your_bot_token_here
   ```

3. **Set the Authentication Passphrase (Optional)**: You can create a custom authentication passphrase by setting the `BOT_AUTH` environment variable.

   ```bash
   export BOT_AUTH=your_bot_auth_here
   ```

   By default "this bot is cool" is used as the passphrase.
   This passphrase will be used to authorize users before they can interact with the bot. Otherwise, you will receive an error message.

4. **Build and Run the Bot**: Build the bot executable with.

   ```bash
   go build ./cmd/bot
   ```

   Then, run the bot:

   ```bash
   ./bot
   ```

   The bot will now be running and ready to accept commands from users.

## Usage

Once the bot is running, users can interact with it through the Telegram chat interface. The bot supports the following commands:

- `/help`: Shows list of commands and what they do.
- `/create`: Create a new task.
- `/tasks`: Interact with existing tasks. From here you can select task and then proceed to complete it, check info, etc.
- `/score`: Shows the score of all participants in 3 periods: current month, last month, and total score. 
- `/monthly`: Returns the list of all actions done during this month.



## Development and Contributions

This is the initial version of the Telegram Track Routines Bot, so there may be room for improvement. If you encounter any issues or have suggestions for new features, please feel free to open an issue or submit a pull request on the project's repository.
Right now bot is using Long polling, in future version it may switch to Webhooks. Also plan to switch to MySQL or PostGreSQL instead of SQLite.
