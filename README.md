[English readme](https://github.com/orrikado/sentimenta/blob/main/README.md) • [Русский readme](https://github.com/orrikado/sentimenta/blob/main/README.ru.md)


# Sentimenta — Mood Tracker with AI Suggestions

**Sentimenta** is a web application for daily mood tracking and analysis. Users rate their mood on a scale from 1 to 5, select emotions from a list or add their own, and optionally write a short note about their day. After saving the entry, the system provides a personalized suggestion generated by AI. The interface is built with SvelteKit and Tailwind CSS, ensuring high performance and responsive design.

## Key Features

![изображение](https://github.com/user-attachments/assets/8a95ff92-5552-46ea-8fe0-b89c64da3ff9)

* **Daily Mood Logging:** Users select a rating from 1 to 5, choose emotions, and optionally write a note about their day.
* **AI Suggestions:** After each entry, Sentimenta generates a recommendation for improving mood.
* **Authentication:** Secure authentication via JWT and OAuth (GitHub/Google).
* **Data Storage:** All entries are stored in a PostgreSQL database.
* **Statistics:** A chart displays mood rating trends over the past month.
* **Monitoring:** Prometheus is used for metrics collection, and Grafana for visualization.
* **Containerization:** The project is fully containerized with Docker (using `docker-compose` and Traefik for routing).

Sentimenta is optimized for different devices. The responsive design makes it easy to log mood entries and view personal stats. Clear data visualizations and a fast interface ensure a smooth and enjoyable user experience.

## Technologies

* **Frontend:** SvelteKit, Tailwind CSS
* **Backend:** Go (Golang), Echo web framework
* **Database:** PostgreSQL
* **Monitoring:** Prometheus, Grafana

## Installation and Launch

1. Clone the repository and navigate into its directory.

2. Copy the `dotenv_template` file to `.env` and configure the required environment variables (database settings, OAuth, secrets, etc.).

3. Make sure Docker and Docker Compose are installed.

4. Launch all services with the following command:

   ```bash
   docker-compose up -d
   ```

5. The application will be available at the configured domain (e.g., `https://sentimenta.example.com`).

## Future Plans

* Release mobile apps for iOS and Android.
* Add reminders and expanded personalization options.

## License

This project is licensed under the MIT License (see the [LICENSE](LICENSE) file for details).
