# quananqr1/Dockerfile
FROM node:21.7.1

# Set Node to use legacy OpenSSL provider
ENV NODE_OPTIONS='--openssl-legacy-provider'

# Set working directory
WORKDIR /app

# Copy package files first to leverage Docker cache
COPY package.json package-lock.json* ./

# Install dependencies
RUN npm install

# Copy the rest of the application
COPY . .

# Expose the port Next.js runs on
EXPOSE 3000

# Set environment variables for development
ENV NODE_ENV=development
ENV NEXT_TELEMETRY_DISABLED=1

# Start the application in development mode
CMD ["npm", "run", "dev"]
# RUN npm run build

# CMD ["npm", "start"]