import { useFocusEffect, useRouter } from 'expo-router';
import { useCallback, useState } from 'react';
import { Button, SafeAreaView, StyleSheet, Text, View } from 'react-native';

import type { Dashboard } from '../src/generated/api';
import { logout } from '../src/domain/auth';
import { getDashboard } from '../src/domain/finance';

export default function DashboardScreen() {
  const router = useRouter();
  const [data, setData] = useState<Dashboard | null>(null);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async () => {
    setError(null);
    try {
      setData(await getDashboard());
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha ao carregar dashboard');
    }
  }, []);

  useFocusEffect(
    useCallback(() => {
      void load();
      return () => undefined;
    }, [load])
  );

  async function signOut() {
    await logout();
    router.replace('/login');
  }

  return (
    <SafeAreaView style={styles.safe}>
      <View style={styles.container}>
        <Text style={styles.title}>Dashboard</Text>
        {error ? <Text style={styles.error}>{error}</Text> : null}
        <Card label="Receita mensal" value={String(data?.monthly_revenue_cents ?? 0)} />
        <Card label="MRR" value={String(data?.mrr_cents ?? 0)} />
        <Card label="Pendências" value={String(data?.pending_payments ?? 0)} />
        <View style={styles.actions}>
          <Button title="Atualizar" onPress={() => void load()} />
          <Button title="Sair" color="#b91c1c" onPress={() => void signOut()} />
        </View>
      </View>
    </SafeAreaView>
  );
}

function Card({ label, value }: { label: string; value: string }) {
  return (
    <View style={styles.card}>
      <Text style={styles.cardLabel}>{label}</Text>
      <Text style={styles.cardValue}>{value}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  safe: { flex: 1, backgroundColor: '#f8fafc' },
  container: { flex: 1, padding: 20, gap: 12 },
  title: { fontSize: 28, fontWeight: '700', color: '#0f172a' },
  error: { color: '#b91c1c' },
  card: {
    borderRadius: 14,
    padding: 14,
    backgroundColor: '#fff',
    borderWidth: 1,
    borderColor: '#e2e8f0',
  },
  cardLabel: { color: '#475569', fontSize: 13 },
  cardValue: { color: '#0f172a', fontWeight: '700', fontSize: 20, marginTop: 4 },
  actions: { marginTop: 16, flexDirection: 'row', gap: 12 },
});
