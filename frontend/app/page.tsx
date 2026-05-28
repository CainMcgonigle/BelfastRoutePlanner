import Map from '../components/Map'
import TripPlanner from '../components/TripPlanner'

export default function Home() {
  return (
    <main
      style={{
        display: 'flex',
        height: '100vh',
        width: '100vw',
        overflow: 'hidden',
      }}
    >
      <div style={{ flex: '0 0 70%', position: 'relative' }}>
        <Map />
      </div>
      <div
        style={{
          flex: '0 0 30%',
          overflowY: 'auto',
          borderLeft: '1px solid #1e293b',
        }}
      >
        <TripPlanner />
      </div>
    </main>
  )
}
