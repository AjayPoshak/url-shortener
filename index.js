

const crypto = require('crypto');

class URLShortener {
    constructor() {
        this.urlDatabase = new Map();
    }

    // Convert number to base62 string
    toBase62(num) {
        const chars = '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';
        let result = '';
        
        do {
            result = chars[num % 62] + result;
            num = Math.floor(num / 62);
        } while (num > 0);
        
        return result;
    }

    // Generate short URL
    shortenURL(longURL) {
        // Create MD5 hash of the long URL
        const hash = crypto.createHash('md5').update(longURL).digest('hex');
        
        // Take first 8 characters of hash and convert to number
        const hashPrefix = hash.substr(0, 8);
        const number = parseInt(hashPrefix, 16);
        
        // Convert to base62 for shorter, URL-safe string
        let shortCode = this.toBase62(number);
        
        // Ensure minimum length of 6 characters by padding if necessary
        while (shortCode.length < 6) {
            shortCode = '0' + shortCode;
        }
        
        // Store in database
        this.urlDatabase.set(shortCode, longURL);
        
        return shortCode;
    }

    // Retrieve original URL
    getLongURL(shortCode) {
        return this.urlDatabase.get(shortCode) || null;
    }

    // Check collision and handle if necessary
    handleCollision(shortCode, longURL) {
        if (this.urlDatabase.has(shortCode)) {
            const existingURL = this.urlDatabase.get(shortCode);
            if (existingURL !== longURL) {
                // If collision detected, append counter to original URL and rehash
                let counter = 1;
                let newLongURL = `${longURL}#${counter}`;
                let newShortCode = this.shortenURL(newLongURL);
                
                while (this.urlDatabase.has(newShortCode)) {
                    counter++;
                    newLongURL = `${longURL}#${counter}`;
                    newShortCode = this.shortenURL(newLongURL);
                }
                
                return newShortCode;
            }
        }
        return shortCode;
    }
}

// Example usage
const shortener = new URLShortener();

// Test the shortener
const longURL = 'https://www.example.com/very/long/path/with/parameters?param1=value1&param2=value2';
const shortCode = shortener.shortenURL(longURL);
console.log(`Short URL code: ${shortCode}`);
console.log(`Original URL: ${shortener.getLongURL(shortCode)}`);
