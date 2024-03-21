-- Patients table
CREATE TABLE Patients (
    PatientID SERIAL PRIMARY KEY,
    Name VARCHAR(100),
    DateOfBirth DATE,
    Gender VARCHAR(10),
    ContactInfo VARCHAR(100)
);

-- Doctors table
CREATE TABLE Doctors (
    DoctorID SERIAL PRIMARY KEY,
    Name VARCHAR(100),
    Specialization VARCHAR(100),
    ContactInfo VARCHAR(100)
);

-- Appointments table
CREATE TABLE Appointments (
    AppointmentID SERIAL PRIMARY KEY,
    PatientID INT REFERENCES Patients(PatientID),
    DoctorID INT REFERENCES Doctors(DoctorID),
    AppointmentDateTime TIMESTAMP,
    Status VARCHAR(20)
);

-- MedicalRecords table
CREATE TABLE MedicalRecords (
    RecordID SERIAL PRIMARY KEY,
    PatientID INT REFERENCES Patients(PatientID),
    DoctorID INT REFERENCES Doctors(DoctorID),
    Date DATE,
    Diagnosis TEXT,
    Treatment TEXT
);

-- Medications table
CREATE TABLE Medications (
    MedicationID SERIAL PRIMARY KEY,
    Name VARCHAR(100),
    Manufacturer VARCHAR(100),
    Dosage VARCHAR(50),
    Price DECIMAL(10, 2)
);

-- Prescriptions table
CREATE TABLE Prescriptions (
    PrescriptionID SERIAL PRIMARY KEY,
    RecordID INT REFERENCES MedicalRecords(RecordID),
    MedicationID INT REFERENCES Medications(MedicationID),
    DosageInstructions TEXT,
    Quantity INT
);

-- Patients table
CREATE TABLE Patients (
    PatientID SERIAL PRIMARY KEY,
    Name VARCHAR(100),
    DateOfBirth DATE,
    Gender VARCHAR(10),
    ContactInfo VARCHAR(100),
    Username VARCHAR(50) UNIQUE, -- Add a username field for online registration
    Password VARCHAR(50) -- Add a password field for online registration
);

-- RegistrationRequests table
CREATE TABLE RegistrationRequests (
    RequestID SERIAL PRIMARY KEY,
    Name VARCHAR(100),
    DateOfBirth DATE,
    Gender VARCHAR(10),
    ContactInfo VARCHAR(100),
    Username VARCHAR(50) UNIQUE,
    Password VARCHAR(50),
    RequestDateTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    Status VARCHAR(20) DEFAULT 'Pending' -- Status can be 'Pending', 'Approved', or 'Rejected'
);

-- ConsultationRequests table
CREATE TABLE ConsultationRequests (
    RequestID SERIAL PRIMARY KEY,
    PatientID INT REFERENCES Patients(PatientID),
    DoctorID INT REFERENCES Doctors(DoctorID),
    RequestDateTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    Status VARCHAR(20) DEFAULT 'Pending' -- Status can be 'Pending', 'Approved', or 'Rejected'
);

-- Consultations table
CREATE TABLE Consultations (
    ConsultationID SERIAL PRIMARY KEY,
    PatientID INT REFERENCES Patients(PatientID),
    DoctorID INT REFERENCES Doctors(DoctorID),
    StartTime TIMESTAMP,
    EndTime TIMESTAMP,
    Status VARCHAR(20) DEFAULT 'Scheduled' -- Status can be 'Scheduled', 'In Progress', 'Completed', or 'Canceled'
);

-- DiagnosisResults table
CREATE TABLE DiagnosisResults (
    ResultID SERIAL PRIMARY KEY,
    PatientID INT REFERENCES Patients(PatientID),
    DoctorID INT REFERENCES Doctors(DoctorID),
    UploadDateTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ResultFile BYTEA, -- Store the diagnosis result file (you may need to adjust the data type based on your database's requirements)
    Notes TEXT
);

-- ChatMessages table
CREATE TABLE ChatMessages (
    MessageID SERIAL PRIMARY KEY,
    SenderID INT,
    ReceiverID INT,
    MessageType VARCHAR(20),
    MessageText TEXT,
    MessageDateTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (SenderID) REFERENCES Patients(PatientID) ON DELETE CASCADE,
    FOREIGN KEY (ReceiverID) REFERENCES Doctors(DoctorID) ON DELETE CASCADE
);

