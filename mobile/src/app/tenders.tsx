import { View, Text, ScrollView, StyleSheet, TouchableOpacity, TextInput } from 'react-native';
import { useColorScheme } from 'react-native';
import { useState } from 'react';
import { Colors } from '@/constants/theme';
import { getTenders, createTender } from '@/lib/api';

export default function TendersScreen() {
  const scheme = useColorScheme();
  const colors = Colors[scheme === 'unspecified' ? 'light' : scheme];
  const [searchQuery, setSearchQuery] = useState('');
  const [showModal, setShowModal] = useState(false);

  return (
    <View style={[styles.container, { backgroundColor: colors.background }]}>
      <ScrollView style={styles.scrollView}>
        <View style={styles.header}>
          <Text style={[styles.title, { color: colors.text }]}>Tenders</Text>
          <Text style={[styles.subtitle, { color: colors.textSecondary }]}>
            Business opportunities and contracts
          </Text>
        </View>

        <View style={styles.searchContainer}>
          <TextInput
            style={[styles.searchInput, { color: colors.text, backgroundColor: colors.card }]}
            placeholder="Search tenders..."
            placeholderTextColor={colors.textSecondary}
            value={searchQuery}
            onChangeText={setSearchQuery}
          />
        </View>

        <View style={styles.content}>
          <View style={styles.headerRow}>
            <Text style={[styles.sectionTitle, { color: colors.text }]}>
              Open Tenders
            </Text>
            <TouchableOpacity
              style={[styles.addButton, { backgroundColor: colors.primary }]}
              onPress={() => setShowModal(true)}
            >
              <Text style={styles.addButtonText}>+ Post Tender</Text>
            </TouchableOpacity>
          </View>
          
          {/* Sample tender cards */}
          {[
            { title: 'School Renovation Project', organization: 'Ministry of Education', deadline: '2024-08-15', budget: 'KES 5M' },
            { title: 'IT Infrastructure Upgrade', organization: 'County Government', deadline: '2024-09-01', budget: 'KES 12M' },
            { title: 'Supply of Office Equipment', organization: 'Public Service Commission', deadline: '2024-07-30', budget: 'KES 2M' },
          ].map((tender, index) => (
            <TouchableOpacity
              key={index}
              style={[styles.card, { backgroundColor: colors.card, borderColor: colors.border }]}
            >
              <View style={styles.cardHeader}>
                <Text style={[styles.tenderTitle, { color: colors.text }]}>{tender.title}</Text>
                <Text style={[styles.budget, { color: colors.primary }]}>{tender.budget}</Text>
              </View>
              <Text style={[styles.organization, { color: colors.textSecondary }]}>{tender.organization}</Text>
              <View style={styles.cardFooter}>
                <Text style={[styles.deadline, { color: colors.textSecondary }]}>
                  Deadline: {tender.deadline}
                </Text>
              </View>
            </TouchableOpacity>
          ))}
        </View>
      </ScrollView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  scrollView: {
    flex: 1,
  },
  header: {
    padding: 20,
    paddingTop: 60,
  },
  title: {
    fontSize: 32,
    fontWeight: 'bold',
  },
  subtitle: {
    fontSize: 14,
    marginTop: 4,
  },
  searchContainer: {
    marginHorizontal: 20,
    marginBottom: 20,
  },
  searchInput: {
    height: 48,
    borderRadius: 12,
    paddingHorizontal: 16,
    fontSize: 14,
  },
  content: {
    padding: 20,
  },
  headerRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 16,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: '600',
  },
  addButton: {
    paddingHorizontal: 16,
    paddingVertical: 8,
    borderRadius: 8,
  },
  addButtonText: {
    color: '#fff',
    fontSize: 12,
    fontWeight: '600',
  },
  card: {
    padding: 16,
    borderRadius: 12,
    marginBottom: 12,
    borderWidth: 1,
  },
  cardHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    marginBottom: 8,
  },
  tenderTitle: {
    fontSize: 16,
    fontWeight: '600',
    flex: 1,
  },
  budget: {
    fontSize: 14,
    fontWeight: '600',
  },
  organization: {
    fontSize: 14,
    marginTop: 4,
  },
  cardFooter: {
    marginTop: 8,
  },
  deadline: {
    fontSize: 12,
  },
});
