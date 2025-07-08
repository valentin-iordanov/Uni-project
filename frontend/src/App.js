import React, { useState, useEffect, createContext, useContext } from 'react';
import { createRoot } from 'react-dom/client';
import PropTypes from 'prop-types';
import {
  FaStar,
  FaPhone,
  FaEnvelope,
  FaFacebook,
  FaInstagram,
  FaBars,
  FaTimes,
} from 'react-icons/fa';

// ============================================================
// ThemeContext: for dynamic accent color or dark mode
// ============================================================
const ThemeContext = createContext();

function ThemeProvider({ children }) {
  const [accent, setAccent] = useState('#A89E8A'); // olive green accent
  const toggleAccent = () =>
    setAccent(accent === '#A89E8A' ? '#C29D5F' : '#A89E8A'); // gold alternative
  return (
    <ThemeContext.Provider value={{ accent, toggleAccent }}>
      {children}
    </ThemeContext.Provider>
  );
}
ThemeProvider.propTypes = { children: PropTypes.node.isRequired };

// ============================================================
// Utility Hook: useViewportWidth
// ============================================================
function useViewportWidth() {
  const [width, setWidth] = useState(window.innerWidth);
  useEffect(() => {
    const handleResize = () => setWidth(window.innerWidth);
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);
  return width;
}

// ============================================================
// Smooth scroll utility
// ============================================================
function smoothScrollTo(id) {
  const element = document.getElementById(id);
  if (element) element.scrollIntoView({ behavior: 'smooth' });
}

// ============================================================
// Navbar Component
// ============================================================
function Navbar() {
  const [scrolled, setScrolled] = useState(false);
  const [mobileOpen, setMobileOpen] = useState(false);
  const width = useViewportWidth();
  const { accent, toggleAccent } = useContext(ThemeContext);

  // Scroll listener
  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 60);
    window.addEventListener('scroll', onScroll);
    return () => window.removeEventListener('scroll', onScroll);
  }, []);

  // Close mobile menu on navigation
  const handleLinkClick = (id) => {
    smoothScrollTo(id);
    if (mobileOpen) setMobileOpen(false);
  };

  // Menu items
  const sections = [
    { id: 'начало', label: 'Начало' },
    { id: 'за-нас', label: 'За нас' },
    { id: 'галерия', label: 'Галерия' },
    { id: 'ревюта', label: 'Ревюта' },
    { id: 'контакт', label: 'Контакт' },
  ];

  return (
    <header
      className={
        `fixed top-0 left-0 w-full z-50 transition-all bg-opacity-80 backdrop-blur-sm ` +
        (scrolled ? 'bg-white shadow-md' : 'bg-transparent')
      }
    >
      <div className="max-w-7xl mx-auto flex items-center justify-between p-4">
        <div className="flex items-center space-x-2">
          <span className="text-4xl font-serif text-gray-800">Gracia House</span>
          {/* Accent toggle for demonstration */}
          <button
            onClick={toggleAccent}
            aria-label="Toggle Accent"
            className="ml-2"
            style={{ color: accent }}
          >
            🌟
          </button>
        </div>

        {/* Desktop Navigation */}
        {width >= 768 ? (
          <nav className="flex items-center space-x-6">
            {sections.map(({ id, label }) => (
              <button
                key={id}
                onClick={() => handleLinkClick(id)}
                className="text-base font-sans text-gray-800 hover:text-opacity-75"
              >
                {label}
              </button>
            ))}
            <button
              onClick={() => handleLinkClick('контакт')}
              className="px-4 py-2 font-sans text-base text-white rounded-md"
              style={{ backgroundColor: accent }}
            >
              Контакт
            </button>
          </nav>
        ) : (
          // Mobile toggle
          <button
            onClick={() => setMobileOpen((o) => !o)}
            aria-label="Toggle Menu"
            className="text-2xl text-gray-800"
          >
            {mobileOpen ? <FaTimes /> : <FaBars />}
          </button>
        )}
      </div>

      {/* Mobile Drawer */}
      {mobileOpen && width < 768 && (
        <div className="bg-white shadow-lg">
          {sections.map(({ id, label }) => (
            <button
              key={id}
              onClick={() => handleLinkClick(id)}
              className="block w-full text-left px-6 py-4 font-sans text-gray-700 border-b"
            >
              {label}
            </button>
          ))}
        </div>
      )}
    </header>
  );
}

// ============================================================
// Hero Component
// ============================================================
const images = [
  'https://images.unsplash.com/photo-1506744038136-46273834b3fb?auto=format&fit=crop&w=1050&q=80',
  'https://images.unsplash.com/photo-1512918728675-ed5a9ecdebfd?auto=format&fit=crop&w=1050&q=80',
  'https://images.unsplash.com/photo-1522708323590-d24dbb6b0267?auto=format&fit=crop&w=1050&q=80',
];

function Hero() {
  const [index, setIndex] = useState(0);
  const { accent } = useContext(ThemeContext);

  useEffect(() => {
    const interval = setInterval(() => {
      setIndex((i) => (i + 1) % images.length);
    }, 5000);
    return () => clearInterval(interval);
  }, []);

  return (
    <section id="начало" className="relative h-screen overflow-hidden">
      {/* Slides */}
      {images.map((src, i) => (
        <img
          key={src}
          src={src}
          alt={`Slide ${i + 1}`}
          className={
            `absolute inset-0 w-full h-full object-cover transition-opacity duration-1500 ` +
            (i === index ? 'opacity-100' : 'opacity-0')
          }
        />
      ))}
      {/* Overlay Text */}
      <div className="absolute inset-0 flex flex-col items-center justify-center bg-black bg-opacity-30 text-center px-4">
        <h1 className="text-4xl md:text-5xl lg:text-6xl font-serif text-white mb-4">
          Добре дошли в нашата уютна къща за гости
        </h1>
        <h2 className="text-2xl md:text-3xl font-serif text-white mb-6">
          Идеалното място за вашата семейна почивка с приятели или тиймбилдинг.
        </h2>
        <button
          onClick={() => smoothScrollTo('за-нас')}
          className="px-6 py-3 text-base font-sans text-white rounded-md"
          style={{ backgroundColor: accent }}
        >
          Разгледай къщата
        </button>
      </div>
    </section>
  );
}

// ============================================================
// About Component
// ============================================================
function About() {
  const amenities = [
    { icon: '🛏️', label: '7 стаи' },
    { icon: '🏊', label: 'Басейн' },
    { icon: '🛁', label: 'Джакузи' },
    { icon: '🔥', label: 'Сауна' },
    { icon: '🌳', label: 'Градина' },
    { icon: '🍖', label: 'Барбекю' },
    { icon: '🚗', label: 'Паркинг' },
    { icon: '🎤', label: 'Караоке' },
    { icon: '📶', label: 'WiFi' },
  ];

  return (
    <section id="за-нас" className="py-20 bg-gray-50">
      <div className="max-w-5xl mx-auto px-4 text-center">
        <h2 className="text-3xl font-serif text-gray-900 mb-6">
          Gracia House - Вашето място за пълноценна почивка и релакс
        </h2>
        <p className="max-w-3xl mx-auto text-base font-sans text-gray-700 mb-12">
          Gracia House включва две самостоятелни сгради, уютно обзаведени и оборудвани с всички
          съвременни удобства. Всяка стая разполага със собствен санитарен възел, климатик и
          безплатен WiFi. На територията ще откриете басейн, джакузи, сауна и просторна градина,
          идеална за почивка на открито.
        </p>
        <div className="grid grid-cols-1 md:grid-cols-6 lg:grid-cols-12 gap-8">
          {amenities.map(({ icon, label }, idx) => (
            <div key={idx} className="flex flex-col items-center">
              <div className="text-4xl mb-2">{icon}</div>
              <span className="text-base font-sans text-gray-800">{label}</span>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

// ============================================================
// Gallery Component
// ============================================================
function Gallery() {
  // create 12 placeholder URLs
  const gallery = Array.from({ length: 12 }, (_, i) =>
    `https://source.unsplash.com/800x600/?guesthouse,interior,${i}`
  );

  const [lightbox, setLightbox] = useState({ isOpen: false, index: 0 });
  const { accent } = useContext(ThemeContext);

  const openBox = (i) => setLightbox({ isOpen: true, index: i });
  const closeBox = () => setLightbox({ isOpen: false, index: 0 });
  const prev = () =>
    setLightbox((lb) => ({ isOpen: true, index: (lb.index + gallery.length - 1) % gallery.length }));
  const next = () =>
    setLightbox((lb) => ({ isOpen: true, index: (lb.index + 1) % gallery.length }));

  return (
    <section id="галерия" className="py-20 bg-white">
      <div className="max-w-6xl mx-auto px-4">
        <h2 className="text-3xl font-serif text-center text-gray-900 mb-8">
          Разгледайте нашата галерия
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-6 lg:grid-cols-12 gap-8">
          {gallery.map((src, i) => (
            <img
              key={src}
              src={src}
              alt={`Gallery ${i}`}
              className="w-full h-48 object-cover cursor-pointer rounded-lg shadow"
              onClick={() => openBox(i)}
            />
          ))}
        </div>
      </div>

      {lightbox.isOpen && (
        <div className="fixed inset-0 bg-black bg-opacity-80 flex items-center justify-center p-4">
          <button
            onClick={closeBox}
            className="absolute top-4 right-4 text-4xl text-white"
          >
            ×
          </button>
          <button
            onClick={prev}
            className="absolute left-4 text-4xl text-white"
          >
            ‹
          </button>
          <img
            src={gallery[lightbox.index]}
            alt="Lightbox"
            className="max-h-full object-contain"
          />
          <button
            onClick={next}
            className="absolute right-4 text-4xl text-white"
          >
            ›
          </button>
        </div>
      )}
    </section>
  );
}

// ============================================================
// Location Component
// ============================================================
function Location() {
  return (
    <section id="местоположение" className="py-20 bg-gray-50">
      <div className="max-w-6xl mx-auto px-4 text-center">
        <h2 className="text-3xl font-serif text-gray-900 mb-6">Къде се намираме?</h2>
        <div className="w-full h-64 mb-4 rounded-lg overflow-hidden shadow">
          <iframe
            title="Gracia House on Map"
            src="https://www.google.com/maps/embed?pb=!1m18..."
            loading="lazy"
            className="w-full h-full border-0"
          ></iframe>
        </div>
        <button
          onClick={() => window.open('https://www.google.com/maps', '_blank')}
          className="px-6 py-3 font-sans text-base text-white rounded-md"
          style={{ backgroundColor: '#A89E8A' }}
        >
          Отвори в Google Maps
        </button>
      </div>
    </section>
  );
}

// ============================================================
// Reviews Component
// ============================================================
function Reviews() {
  const reviews = [
    {
      stars: 5,
      text: 'Прекрасно място и невероятно обслужване!',
      author: 'Мария Петрова',
      date: 'Юли 2025',
    },
    {
      stars: 4,
      text: 'Уютна атмосфера, ще се върна отново.',
      author: 'Иван Димитров',
      date: 'Юни 2025',
    },
    {
      stars: 5,
      text: 'Препоръчвам за семейства и приятели.',
      author: 'Елена Георгиева',
      date: 'Май 2025',
    },
  ];

  return (
    <section id="ревюта" className="py-20 bg-white">
      <div className="max-w-5xl mx-auto px-4 text-center">
        <h2 className="text-3xl font-serif text-gray-900 mb-8">
          Какво казват нашите гости
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-6 lg:grid-cols-12 gap-8">
          {reviews.map((r, idx) => (
            <div key={idx} className="p-6 border rounded-lg shadow-sm">
              <div className="flex justify-center mb-3">
                {Array.from({ length: r.stars }).map((_, i) => (
                  <FaStar key={i} className="text-yellow-500 mx-1" />
                ))}
              </div>
              <p className="text-base text-gray-700 mb-4">"{r.text}"</p>
              <p className="text-sm font-sans text-gray-500 italic">
                {r.author}, {r.date}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

// ============================================================
// Contact Component
// ============================================================
function Contact() {
  const [form, setForm] = useState({ name: '', email: '', phone: '', message: '' });
  const handleChange = (e) => setForm({ ...form, [e.target.name]: e.target.value });
  const handleSubmit = (e) => {
    e.preventDefault();
    alert(`Thanks, ${form.name}! We'll be in touch.`);
    setForm({ name: '', email: '', phone: '', message: '' });
  };

  const contactInfo = [
    { icon: <FaPhone />, text: '+359 123 456 789' },
    { icon: <FaEnvelope />, text: 'info@graciahouse.bg' },
    { icon: <FaFacebook />, text: 'facebook.com/graciahouse' },
    { icon: <FaInstagram />, text: '@graciahouse' },
  ];

  return (
    <section id="контакт" className="py-20 bg-gray-50">
      <div className="max-w-6xl mx-auto px-4 grid grid-cols-1 md:grid-cols-6 lg:grid-cols-12 gap-8">
        <div className="md:col-span-3 lg:col-span-6">
          <h2 className="text-3xl font-serif text-gray-900 mb-6">Свържете се с нас</h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <input
              type="text"
              name="name"
              value={form.name}
              onChange={handleChange}
              placeholder="Вашето име"
              required
              className="w-full p-2 border rounded"
            />
            <input
              type="email"
              name="email"
              value={form.email}
              onChange={handleChange}
              placeholder="Имейл"
              required
              className="w-full p-2 border rounded"
            />
            <input
              type="tel"
              name="phone"
              value={form.phone}
              onChange={handleChange}
              placeholder="Телефон"
              className="w-full p-2 border rounded"
            />
            <textarea
              name="message"
              value={form.message}
              onChange={handleChange}
              placeholder="Съобщение"
              rows="4"
              className="w-full p-2 border rounded"
            ></textarea>
            <button
              type="submit"
              className="px-6 py-3 text-white rounded-md"
              style={{ backgroundColor: '#A89E8A' }}
            >
              Изпрати
            </button>
          </form>
        </div>
        <div className="md:col-span-3 lg:col-span-6 flex flex-col space-y-4">
          {contactInfo.map(({ icon, text }, idx) => (
            <div key={idx} className="flex items-center space-x-2 text-gray-700">
              {icon}
              <span>{text}</span>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

// ============================================================
// App Component
// ============================================================
export default function App() {
  return (
    <ThemeProvider>
      <Navbar />
      <main className="pt-16">
        <Hero />
        <About />
        <Gallery />
        <Location />
        <Reviews />
        <Contact />
      </main>
    </ThemeProvider>
  );
}
