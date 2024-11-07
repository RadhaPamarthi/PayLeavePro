import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';

function LoginPage() {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [message, setMessage] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const navigate = useNavigate();

    const handleLogin = async (e) => {
        e.preventDefault();
        setIsLoading(true);
        setMessage('');

        try {
            // Add headers to handle CORS and content type
            const config = {
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json'
                }
            };

            const response = await axios.post('http://localhost:8080/login', {
                username,
                password
            }, config);

            console.log('Full Response:', response); // Debug log

            // Check if response has data
            if (response.data) {
                // Check for success status
                if (response.status === 200) {
                    // Store any tokens or user data if needed
                    if (response.data.token) {
                        localStorage.setItem('token', response.data.token);
                    }
                    
                    setMessage('Login successful!');
                    // Short delay before navigation
                    setTimeout(() => {
                        navigate('/dashboard');
                    }, 500);
                } else {
                    setMessage(response.data.message || 'Something went wrong');
                }
            }
        } catch (error) {
            console.error('Login Error:', error);
            
            if (error.response) {
                // The request was made and the server responded with a status code
                // that falls out of the range of 2xx
                console.log('Error Response Data:', error.response.data);
                console.log('Error Response Status:', error.response.status);
                console.log('Error Response Headers:', error.response.headers);
                
                setMessage(error.response.data.message || 'Server error occurred');
            } else if (error.request) {
                // The request was made but no response was received
                console.log('Error Request:', error.request);
                setMessage('No response from server');
            } else {
                // Something happened in setting up the request that triggered an Error
                console.log('Error Message:', error.message);
                setMessage('Failed to send request');
            }
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div style={styles.container}>
            <div style={styles.card}>
                <h2 style={styles.title}>Admin Login</h2>
                <form onSubmit={handleLogin}>
                    <div style={styles.inputGroup}>
                        <label style={styles.label}>Username:</label>
                        <input
                            type="text"
                            value={username}
                            onChange={(e) => setUsername(e.target.value)}
                            style={styles.input}
                            placeholder="Enter your username"
                            disabled={isLoading}
                        />
                    </div>
                    <div style={styles.inputGroup}>
                        <label style={styles.label}>Password:</label>
                        <input
                            type="password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            style={styles.input}
                            placeholder="Enter your password"
                            disabled={isLoading}
                        />
                    </div>
                    <button 
                        type="submit" 
                        style={{
                            ...styles.button,
                            ...(isLoading && styles.buttonDisabled)
                        }}
                        disabled={isLoading}
                    >
                        {isLoading ? 'Logging in...' : 'Login'}
                    </button>
                </form>
                {message && (
                    <p style={{
                        ...styles.message,
                        color: message.includes('successful') ? '#28a745' : '#dc3545'
                    }}>
                        {message}
                    </p>
                )}
            </div>
        </div>
    );
}

const styles = {
    container: {
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        height: '100vh',
        backgroundColor: '#f5f5f5',
    },
    card: {
        width: '100%',
        maxWidth: '400px',
        padding: '30px',
        borderRadius: '12px',
        boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1), 0 1px 3px rgba(0, 0, 0, 0.08)',
        backgroundColor: '#ffffff',
    },
    title: {
        fontSize: '28px',
        fontWeight: '600',
        textAlign: 'center',
        marginBottom: '30px',
        color: '#1a1a1a',
    },
    inputGroup: {
        marginBottom: '20px',
    },
    label: {
        fontSize: '14px',
        fontWeight: '500',
        color: '#4a5568',
        marginBottom: '8px',
        display: 'block',
    },
    input: {
        width: '100%',
        padding: '12px 16px',
        fontSize: '16px',
        borderRadius: '8px',
        border: '1px solid #e2e8f0',
        outline: 'none',
        boxSizing: 'border-box',
        transition: 'border-color 0.2s, box-shadow 0.2s',
        ':focus': {
            borderColor: '#3182ce',
            boxShadow: '0 0 0 3px rgba(66, 153, 225, 0.5)',
        }
    },
    button: {
        width: '100%',
        padding: '14px',
        fontSize: '16px',
        fontWeight: '600',
        borderRadius: '8px',
        border: 'none',
        color: '#fff',
        backgroundColor: '#3182ce',
        cursor: 'pointer',
        marginTop: '20px',
        transition: 'all 0.2s ease',
        ':hover': {
            backgroundColor: '#2c5282',
        }
    },
    buttonDisabled: {
        backgroundColor: '#90cdf4',
        cursor: 'not-allowed',
    },
    message: {
        textAlign: 'center',
        marginTop: '20px',
        padding: '10px',
        borderRadius: '6px',
        fontSize: '14px',
        fontWeight: '500',
    },
};

export default LoginPage;