import React, { useState } from 'react';
import AddEmployee from './AddEmployee'; // Adjust the import path if necessary
import EmployeeList from './EmployeeList'; // Adjust the path if necessary

// Icon Components
const IconUser = () => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="w-4 h-4">
    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
    <circle cx="12" cy="7" r="4" />
  </svg>
);

const IconHome = () => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="w-4 h-4">
    <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
    <polyline points="9 22 9 12 15 12 15 22" />
  </svg>
);

const IconCalendar = () => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="w-4 h-4">
    <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
    <line x1="16" y1="2" x2="16" y2="6" />
    <line x1="8" y1="2" x2="8" y2="6" />
    <line x1="3" y1="10" x2="21" y2="10" />
  </svg>
);

const IconDollar = () => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="w-4 h-4">
    <line x1="12" y1="1" x2="12" y2="23" />
    <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6" />
  </svg>
);

const IconBell = () => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="w-6 h-6">
    <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9" />
    <path d="M13.73 21a2 2 0 0 1-3.46 0" />
  </svg>
);

const Dashboard = () => {
  const [activeTab, setActiveTab] = useState('dashboard');
  const [showNotifications, setShowNotifications] = useState(false);

  const stats = {
    totalEmployees: 156,
    pendingLeaves: 8,
    payrollDue: "Nov 30, 2024",
    todayAttendance: "92%"
  };

  const notifications = [
    { id: 1, text: "New leave request from John Doe", time: "5m ago" },
    { id: 2, text: "Payroll processing deadline tomorrow", time: "1h ago" }
  ];

  return (
    <div className="flex min-h-screen bg-gray-100">
      {/* Sidebar */}
      <div className="w-64 bg-white shadow-lg">
        <div className="p-4">
          <h2 className="text-xl font-bold mb-4">HR Dashboard</h2>
          <nav className="space-y-2">
            {[
              { id: 'dashboard', icon: <IconHome />, label: 'Dashboard' },
              { id: 'addEmployee', icon: <IconUser />, label: 'Add Employee' },
              { id: 'employeeList', icon: <IconUser />, label: 'Employee List' },
              { id: 'leaves', icon: <IconCalendar />, label: 'Leave Requests' },
              { id: 'payroll', icon: <IconDollar />, label: 'Payroll' },
            ].map(item => (
              <button
                key={item.id}
                onClick={() => setActiveTab(item.id)}
                className={`flex items-center space-x-2 p-2 rounded w-full text-left transition-colors duration-200 ${
                  activeTab === item.id ? 'bg-blue-50 text-blue-600' : 'hover:bg-gray-50'
                }`}
              >
                {item.icon}
                <span>{item.label}</span>
              </button>
            ))}
          </nav>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 p-8">
        {/* Header */}
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-2xl font-bold">Welcome back, Admin!</h1>
            <p className="text-gray-600">
              {new Date().toLocaleDateString('en-US', { 
                weekday: 'long', 
                year: 'numeric', 
                month: 'long', 
                day: 'numeric' 
              })}
            </p>
          </div>
          <div className="flex items-center space-x-4 relative">
            <button 
              className="p-2 relative"
              onClick={() => setShowNotifications(!showNotifications)}
            >
              <IconBell />
              <span className="absolute top-0 right-0 w-2 h-2 bg-red-500 rounded-full"></span>
            </button>
            
            {showNotifications && (
              <div className="absolute top-12 right-0 w-80 bg-white rounded-lg shadow-lg p-4 z-10">
                <h3 className="text-sm font-semibold mb-2">Notifications</h3>
                <div className="space-y-2">
                  {notifications.map(notif => (
                    <div key={notif.id} className="flex items-start space-x-2 p-2 hover:bg-gray-50 rounded">
                      <div className="w-2 h-2 bg-blue-500 rounded-full mt-2"></div>
                      <div>
                        <p className="text-sm">{notif.text}</p>
                        <p className="text-xs text-gray-500">{notif.time}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Tabs Content */}
        {activeTab === 'dashboard' && (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
            {[
              { title: 'Total Employees', value: stats.totalEmployees, icon: <IconUser /> },
              { title: 'Active Today', value: stats.todayAttendance, icon: <IconUser /> },
              { title: 'Pending Leaves', value: stats.pendingLeaves, icon: <IconCalendar /> }
            ].map((stat, idx) => (
              <div key={idx} className="bg-white p-6 rounded-lg shadow hover:shadow-md transition-shadow duration-200">
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-sm font-medium text-gray-600">{stat.title}</h3>
                  {stat.icon}
                </div>
                <div className="text-2xl font-bold mb-2">{stat.value}</div>
              </div>
            ))}
          </div>
        )}

        {activeTab === 'addEmployee' && (
          <AddEmployee />
        )}
        
        {activeTab === 'employeeList' && (
          <EmployeeList />
        )}

        {activeTab === 'leaves' && (
          <div>
            <h2 className="text-xl font-bold mb-4">Leave Requests</h2>
            {/* Render leave requests */}
          </div>
        )}

        {activeTab === 'payroll' && (
          <div>
            <h2 className="text-xl font-bold mb-4">Payroll</h2>
            {/* Render payroll details */}
          </div>
        )}
      </div>
    </div>
  );
};

export default Dashboard;
